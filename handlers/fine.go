package handlers

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
)

var session *gocql.Session

type Fine struct {
    ID     gocql.UUID `json:"id"`
    Player string     `json:"player"`
    Reason string     `json:"reason"`
    Amount float64    `json:"amount"`
    Date   time.Time  `json:"date"`
}

func init() {
    cluster := gocql.NewCluster("cassandra")
    cluster.Keyspace = "league"
    cluster.Consistency = gocql.Quorum
    var err error
    session, err = cluster.CreateSession()
    if err != nil {
        panic(err)
    }
}

func CreateFine(c *gin.Context) {
    var fine Fine
    if err := c.ShouldBindJSON(&fine); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fine.ID = gocql.TimeUUID()
    fine.Date = time.Now()

    if err := session.Query(`INSERT INTO fines (id, player, reason, amount, date) VALUES (?, ?, ?, ?, ?)`,
        fine.ID, fine.Player, fine.Reason, fine.Amount, fine.Date).Exec(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save fine"})
        return
    }

    fmt.Printf("EVENT: fine.issued -> %+v\n", fine)
    c.JSON(http.StatusCreated, fine)
}

func ListFines(c *gin.Context) {
    var fines []Fine
    iter := session.Query("SELECT id, player, reason, amount, date FROM fines").Iter()
    var fine Fine
    for iter.Scan(&fine.ID, &fine.Player, &fine.Reason, &fine.Amount, &fine.Date) {
        fines = append(fines, fine)
    }
    c.JSON(http.StatusOK, fines)
}
