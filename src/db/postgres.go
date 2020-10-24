package db

import (
    "time"

    "worko.tech/iam/src/models"
    "worko.tech/iam/src/utils"

    "github.com/jinzhu/gorm"
    "github.com/rs/zerolog/log"
)

func ConnectToPostgres() (*gorm.DB) {
    // Database connection, dao and handlers
    dbConnected := false
    var postgres *gorm.DB
    var err error

    for dbConnected != true {
        postgres, err = newPostgresConn()
        if err != nil {
            log.Warn().Msgf("Unable to connect to postgres (%v). Retrying...", err.Error())
            time.Sleep(2000 * time.Millisecond)
        } else {
            dbConnected = true
        }
    }

    log.Info().Msg("Successfully connect to postgres")
    return postgres
}

func newPostgresConn() (*gorm.DB, error) {
    hostString := "host=" + utils.GetEnv("POSTGRES_HOST", "localhost")
    portString := "port=" + utils.GetEnv("POSTGRES_PORT", "5432")
    dbString := "dbname=" + utils.GetEnv("POSTGRES_DATABASE", "iam")
    pwdString := "password=" + utils.GetEnv("POSTGRES_PASSWORD", "postgres")
    userString := "user=" + utils.GetEnv("POSTGRES_USER",  "postgres")

    connectionString := hostString + " " + portString + " " + dbString + " " + userString + " " + pwdString + " sslmode=disable"
    db, err := gorm.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    db.AutoMigrate(&models.User{})
    return db, err
}
