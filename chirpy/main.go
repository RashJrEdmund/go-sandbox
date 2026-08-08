/*
SEE WRITING struct of a request handler here
	https://pkg.go.dev/net/http#ResponseWriter.WriteHeader
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync/atomic"
	"time"

	"github.com/RashJrEdmund/go-sandbox/chirpy/internal/auth"
	"github.com/RashJrEdmund/go-sandbox/chirpy/internal/database"
	"github.com/RashJrEdmund/go-sandbox/chirpy/internal/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// MIDDLEWARES

type apiConfig struct {
	// The atomic.Int32 type is a really cool standard-library type that allows us to safely increment and read
	// an integer value across multiple goroutines (HTTP requests). https://pkg.go.dev/sync/atomic#Int32
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incrementing metrics")
		apiCfg.fileServerHits.Add(1)

		fmt.Println(apiCfg.fileServerHits.Load())
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	restTemplate := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, apiCfg.fileServerHits.Load())

	utils.RespondWithPlainText(w, http.StatusOK, restTemplate)
}

func (apiCfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileServerHits.Store(0)
	resText := fmt.Sprintf("Hits: %d", apiCfg.fileServerHits.Load())

	if apiCfg.platform != "dev" {
		utils.RespondWithPlainText(w, http.StatusForbidden, "Forbidden")
		return
	}

	err := apiCfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to reset users.")
		return
	}

	utils.RespondWithPlainText(w, http.StatusOK, resText)
}

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var reqData utils.CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid JSON.")
		return
	}

	hashedPassword, err := auth.HashPassword(reqData.Password)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to hash password.")
		return
	}

	newUser, err := apiCfg.dbQueries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          reqData.Email,
			HashedPassword: hashedPassword,
		},
	)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to create user.")
		return
	}

	userResp := utils.UserResponse{
		ID:          newUser.ID.String(),
		Email:       newUser.Email,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		IsChirpyRed: newUser.IsChirpyRed,
	}

	utils.RespondWithJSON(w, http.StatusCreated, userResp)
}

func (apiCfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var reqData utils.UpdateUserRequest

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	userIDFromToken, err := auth.ValidateJWT(refreshToken, apiCfg.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid JSON.")
		return
	}

	user, err := apiCfg.dbQueries.GetUserByID(r.Context(), userIDFromToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	// Creating new hash

	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to hash password.")
		return
	}

	err = apiCfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             user.ID,
		Email:          reqData.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to update user.")
		return
	}

	updatedAt := time.Now().UTC()

	userResp := utils.UserResponse{
		ID:          user.ID.String(),
		Email:       reqData.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   updatedAt,
		IsChirpyRed: user.IsChirpyRed,
	}

	utils.RespondWithJSON(w, http.StatusOK, userResp)
}

func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	var reqData utils.LoginRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid JSON.")
		return
	}

	user, err := apiCfg.dbQueries.GetUserByEmail(r.Context(), reqData.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password.")
		return
	}

	isValid, err := auth.CheckPasswordHash(reqData.Password, user.HashedPassword)
	if !isValid || err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password.")
		return
	}

	token, err := auth.MakeJWT(
		user.ID,
		apiCfg.jwtSecret,
	)

	refreshToken, err := apiCfg.dbQueries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:  auth.MakeRefreshToken(),
			UserID: user.ID,
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to create refresh token.")
		return
	}

	userResp := utils.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken.Token,
		UserResponse: utils.UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
	}

	utils.RespondWithJSON(w, http.StatusOK, userResp)
}

func (apiCfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	refreshTokenDB, err := apiCfg.dbQueries.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	if refreshTokenDB.RevokedAt.Valid {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token has been revoked.")
		return
	}

	user, err := apiCfg.dbQueries.GetUserByID(r.Context(), refreshTokenDB.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	// create new access token with jwt.

	newAccessToken, err := auth.MakeJWT(
		user.ID,
		apiCfg.jwtSecret,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to create access token.")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.RefreshTokenResponse{Token: newAccessToken})
	return
}

func (apiCfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	err = apiCfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to revoke token.")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, utils.RefreshTokenResponse{Token: "Token revoked"})
	return
}

func (apiCfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var reqData utils.CreateChirpRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid JSON.")
		return
	}

	userIDFromToken, err := auth.GetUserIdFromBearerTokenHeader(r.Header, apiCfg.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	if len(reqData.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	reqData.Body = utils.RemoveProfanity(reqData.Body)

	newChirp, err := apiCfg.dbQueries.CreateChirp(r.Context(),
		database.CreateChirpParams{
			Body:   reqData.Body,
			UserID: userIDFromToken,
		})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to create chirp.")
		return
	}

	chirpResp := utils.CreateChirpResponse{
		ID:        newChirp.ID.String(),
		Body:      newChirp.Body,
		UserId:    newChirp.UserID.String(),
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusCreated, chirpResp)
}

func (apiCfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		authorUUID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid author_id.")
			return
		}
		chirps, err = apiCfg.dbQueries.ListChirpsByAuthorID(r.Context(), authorUUID)
	} else {
		chirps, err = apiCfg.dbQueries.ListAllChirps(r.Context())
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to get chirps.")
		return
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	chirpResp := []utils.CreateChirpResponse{}
	for _, chirp := range chirps {
		chirpResp = append(chirpResp, utils.CreateChirpResponse{
			ID:        chirp.ID.String(),
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, chirpResp)
}

func (apiCfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId") // see: https://pkg.go.dev/net/http#Request.PathValue

	chirpIdUUID, err := uuid.Parse(chirpId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid chirp_id.")
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirpByID(r.Context(), chirpIdUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Failed to get chirp.")
		return
	}

	chirpResp := utils.CreateChirpResponse{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, chirpResp)
}

func (apiCfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId") // see: https://pkg.go.dev/net/http#Request.PathValue

	userIDFromToken, err := auth.GetUserIdFromBearerTokenHeader(r.Header, apiCfg.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token.")
		return
	}

	chirpIdUUID, err := uuid.Parse(chirpId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid chirp_id.")
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirpByID(r.Context(), chirpIdUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Failed to get chirp.")
		return
	}

	if chirp.UserID != userIDFromToken {
		utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp.")
		return
	}

	err = apiCfg.dbQueries.DeleteChirpById(r.Context(), chirpIdUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to delete chirp.")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (apiCfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != apiCfg.polkaKey {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var reqData utils.PolkaWebhookRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid JSON.")
		return
	}

	if reqData.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(reqData.Data.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong. Invalid user_id.")
		return
	}

	_, err = apiCfg.dbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found.")
		return
	}

	err = apiCfg.dbQueries.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Failed to upgrade user.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// NON-API_CONFIG METHOD ROUTE HANDLERS

// ------------------Health Handler----------------------------------------

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK")) // w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	DB, err := sql.Open("postgres", dbURL)
	PLATFORM := os.Getenv("PLATFORM")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	POLKA_KEY := os.Getenv("POLKA_KEY")

	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	mu := http.NewServeMux()

	const PORT = "8080"

	const rootDir = "."

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mu,
	}

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      database.New(DB),
		platform:       PLATFORM,
		jwtSecret:      JWT_SECRET,
		polkaKey:       POLKA_KEY,
	}

	/*
		Now that the path is no longer "/", we need to fix this by using http.StripPrefix
			Read here: https://pkg.go.dev/net/http#StripPrefix
	*/
	mu.Handle("/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/", http.FileServer(http.Dir(rootDir))),
		),
	)

	mu.Handle("GET /admin/metrics/", http.HandlerFunc(apiCfg.metricsHandler))

	mu.Handle("POST /admin/reset/", http.HandlerFunc(apiCfg.resetHandler))

	mu.HandleFunc("GET /api/healthz/", healthzHandler)

	mu.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mu.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)

	mu.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mu.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mu.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)

	mu.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mu.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mu.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.getChirpByIDHandler)
	mu.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.deleteChirpByIDHandler)

	mu.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebhookHandler)

	fmt.Printf("Serving files from %s on port %s\n", rootDir, PORT)
	log.Fatal(server.ListenAndServe())
}
