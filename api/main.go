package main

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
import (
	"api/common"
	"api/handlers"
	"api/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const defaultPprofPort = "8080"
const defaultAPIPort = "80"
const defaultLogPath = "./logs/service.log"

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func initPprof(pprofPort string) {
	log, err := common.GetLogger()
	if err != nil {
		fmt.Println(err)
	}
	addr := fmt.Sprintf("localhost: %v", pprofPort)
	go func() {
		log.Infoln("Starting pprof on %v", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Errorf("failed running pprof: %v", err)
		}
	}()
}

func initEngine() *gin.Engine {

	router := gin.Default()
	router.Use(middleware.LoggerMiddleware)
	router.Use(gin.Recovery())
	router.Use(middleware.AuthMiddleware)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", handlers.SignUpHandler)
		authGroup.POST("/signin", handlers.SignInHandler)
	}

	quizzesGroup := router.Group("/quizzes")
	{
		quizzesGroup.POST("/", handlers.CreateQuizHandler)
		quizzesGroup.GET("/:id", handlers.QuizHandler)
		quizzesGroup.PATCH("/", handlers.UpdateQuizHandler)
		quizzesGroup.POST("/:id/questions", handlers.CreateQuestionHandler)
		quizzesGroup.GET("/:id/questions/:question_id", handlers.QuestionHandler)
		quizzesGroup.POST("/:id/questions/:question_id", handlers.FulfillQuestionHandler)
		quizzesGroup.PUT("/:id/questions/:question_id", handlers.UpdateQuestionHandler)
		quizzesGroup.POST("/:id/settings", handlers.UpdateQuizSettingsHandler)
		quizzesGroup.GET("/:id/settings", handlers.QuizSettingsHandler)
	}

	return router
}

func main() {
	err := common.InitLoggerConfig(&common.Config{Filepath: getEnv("LOG_PATH", defaultLogPath)})
	if err != nil {
		fmt.Println(err)
	}
	log, err := common.GetLogger()
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToLower(getEnv("DEBUG", "")) == "true" {
		log.Level = logrus.DebugLevel
	}

	if strings.ToLower(getEnv("PROFILING_DISABLED", "FALSE")) == "false" {
		pprofPort := getEnv("PPROF_PORT", defaultPprofPort)
		initPprof(pprofPort)
	}
	port := getEnv("API_PORT", defaultAPIPort)
	r := initEngine()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

}
