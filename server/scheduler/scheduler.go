package scheduler

import (
	"context"
	"fmt"
	"log"
	"server/notifier"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var WeeklyProducts []map[string]interface{}

// func StartProductCheckScheduler(pool *pgxpool.Pool, tzName string) {
// 	loc, err := time.LoadLocation(tzName)
// 	if err != nil {
// 		log.Fatal("Invalid timezone:", err)
// 	}

// 	c := cron.New(cron.WithLocation(loc))

// 	// Every Sunday at 11:59 PM in given timezone
// 	_, err = c.AddFunc("59 23 * * 0", func() {
// 		query := `
// 		SELECT id, name, description, price, quantity, created_at
// 		FROM get_weekly_products($1)
// 	`
// 		ctx := context.Background()
// 		rows, err := pool.Query(ctx, query, tzName)
// 		if err != nil {
// 			log.Println("❌ Error fetching weekly products:", err)
// 			return
// 		}
// 		defer rows.Close()

// 		var results []map[string]interface{}
// 		for rows.Next() {
// 			var (
// 				id          int
// 				name        string
// 				description string
// 				price       float64
// 				quantity    int
// 				createdAt   time.Time
// 			)

// 			if err := rows.Scan(&id, &name, &description, &price, &quantity, &createdAt); err != nil {
// 				log.Println("❌ Error scanning product row:", err)
// 				continue
// 			}

// 			results = append(results, map[string]interface{}{
// 				"id":          id,
// 				"name":        name,
// 				"description": description,
// 				"price":       price,
// 				"quantity":    quantity,
// 				"created_at":  createdAt,
// 			})
// 		}

// 		if err := rows.Err(); err != nil {
// 			log.Println("❌ Row iteration error:", err)
// 			return
// 		}

// 		WeeklyProducts = results
// 		log.Printf("✅ Weekly product scan complete. %d product(s) found.\n", len(results))
// 	})

// 	if err != nil {
// 		log.Fatal("Failed to schedule job:", err)
// 	}
// 	c.Start()
// }

// func StartProductCheckScheduler(pool *pgxpool.Pool, tzName string) {

// 	loc, err := time.LoadLocation(tzName)
// 	if err != nil {
// 		log.Fatal("Invalid timezone:", err)
// 	}

// 	c := cron.New(cron.WithLocation(loc))

// 	// ✅ Every Monday at 9:00 AM
// 	_, err = c.AddFunc("13 10 2", func() {
// 		query := `
// 			SELECT id, name, description, price, quantity, created_at
// 			FROM get_weekly_products($1)
// 		`
// 		ctx := context.Background()
// 		rows, err := pool.Query(ctx, query, tzName)
// 		if err != nil {
// 			log.Println("❌ Error fetching weekly products:", err)
// 			return
// 		}
// 		defer rows.Close()

// 		var results []map[string]interface{}
// 		for rows.Next() {
// 			var (
// 				id          int
// 				name        string
// 				description string
// 				price       float64
// 				quantity    int
// 				createdAt   time.Time
// 			)

// 			if err := rows.Scan(&id, &name, &description, &price, &quantity, &createdAt); err != nil {
// 				log.Println("❌ Error scanning product row:", err)
// 				continue
// 			}

// 			results = append(results, map[string]interface{}{
// 				"id":          id,
// 				"name":        name,
// 				"description": description,
// 				"price":       price,
// 				"quantity":    quantity,
// 				"created_at":  createdAt,
// 			})
// 		}

// 		if err := rows.Err(); err != nil {
// 			log.Println("❌ Row iteration error:", err)
// 			return
// 		}

// 		WeeklyProducts = results
// 		log.Printf("✅ Weekly product scan complete. %d product(s) found.\n", len(results))

// 		// 🔔 Insert into notifications table if products were found
// 		if len(results) > 0 {
// 			notifText := fmt.Sprintf("%d products were added this week", len(results))
// 			notifID := uuid.New().String()
// 			expiry := time.Now().Add(15 * 24 * time.Hour) // 15 days expiry

// 			insertNotif := `
// 				INSERT INTO notifications (notif_id, notif_description, notif_count, expiry_time)
// 				VALUES ($1, $2, $3, $4)
// 			`

// 			_, err = pool.Exec(ctx, insertNotif, notifID, notifText, len(results), expiry)
// 			if err != nil {
// 				log.Println("❌ Error inserting notification:", err)
// 			} else {
// 				log.Println("🔔 Notification saved:", notifText)
// 			}
// 		}
// 	})

// 	if err != nil {
// 		log.Fatal("Failed to schedule job:", err)
// 	}

// 	c.Start()
// }Working function

func StartProductCheckScheduler(pool *pgxpool.Pool, tzName string) {
	delay := 1 * time.Minute
	log.Println("⏳ Waiting 1 minutes to run product check...")

	time.AfterFunc(delay, func() {
		query := `
				SELECT id, name, description, price, quantity, created_at
				FROM get_weekly_products($1)
			`
		ctx := context.Background()
		rows, err := pool.Query(ctx, query, tzName)
		if err != nil {
			log.Println("❌ Error fetching weekly products:", err)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var (
				id          int
				name        string
				description string
				price       float64
				quantity    int
				createdAt   time.Time
			)
			if err := rows.Scan(&id, &name, &description, &price, &quantity, &createdAt); err != nil {
				log.Println("❌ Error scanning product row:", err)
				continue
			}
			results = append(results, map[string]interface{}{
				"id":          id,
				"name":        name,
				"description": description,
				"price":       price,
				"quantity":    quantity,
				"created_at":  createdAt,
			})
		}

		if err := rows.Err(); err != nil {
			log.Println("❌ Row iteration error:", err)
			return
		}

		WeeklyProducts = results
		log.Printf("✅ Weekly product scan complete. %d product(s) found.\n", len(results))

		notifText := fmt.Sprintf("%d products were added this week.", len(results))
		expiry := time.Now().AddDate(0, 0, 15) // 15 days from now
		notifID := uuid.New().String()
		insertQuery := `
				INSERT INTO notifications (notif_description, notif_count, expiry_time,notif_id)
				VALUES ($1, $2, $3,$4)
			`
		_, err = pool.Exec(ctx, insertQuery, notifText, len(results), expiry, notifID)
		if err != nil {
			log.Println("❌ Error inserting notification:", err)
		} else {
			log.Println("✅ Notification inserted successfully.")
		}

		if len(WeeklyProducts) > 0 {
			notification := map[string]interface{}{
				"description": fmt.Sprintf("%d products were added this week", len(WeeklyProducts)),
				"count":       len(WeeklyProducts),
				"timestamp":   time.Now().Format(time.RFC3339),
			}

			if err := notifier.SendKafkaNotification("product-notifications", notification); err != nil {
				log.Println("❌ Failed to send Kafka notification:", err)
			}

		}
	})
}
