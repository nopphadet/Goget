package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// โครงสร้างข้อมูลสำหรับแต่ละตาราง
type Product struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Quantity    int    `json:"quantity"`
	Barcode     string `json:"barcode"`
	StockStatus string `json:"stock_status"`
	ImagePath   string `json:"image_path"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProductChange struct {
	ChangesID   int    `json:"changes_id"`
	ChangeType  string `json:"change_type"`
	OldQuantity *int   `json:"old_quantity"` // ใช้ pointer เพราะอาจเป็น NULL
	NewQuantity *int   `json:"new_quantity"` // ใช้ pointer เพราะอาจเป็น NULL
	ChangedAt   string `json:"changed_at"`
	ChangedBy   string `json:"changed_by"`
	CreatedAt   string `json:"created_at"`
	ProductID   *int   `json:"product_id"` // ใช้ pointer เพราะอาจเป็น NULL
}

type Reservation struct {
	ReserveID        int     `json:"reserve_id"`
	UserID           int     `json:"user_id"`
	ProductID        int     `json:"product_id"`
	Quantity         int     `json:"quantity"`
	Status           string  `json:"status"`
	ExpiresAt        *string `json:"expires_at"` // ใช้ pointer เพราะอาจเป็น NULL
	CreatedAt        string  `json:"created_at"`
	ActualQuantity   int     `json:"actual_quantity"`
	ReturnedQuantity int     `json:"returned_quantity"`
	UpdatedAt        *string `json:"updated_at"` // ใช้ pointer เพราะอาจเป็น NULL
}

type Return struct {
	ReturnID   int    `json:"return_id"`
	ReserveID  *int   `json:"reserve_id"` // ใช้ pointer เพราะอาจเป็น NULL
	UserID     *int   `json:"user_id"`    // ใช้ pointer เพราะอาจเป็น NULL
	ProductID  *int   `json:"product_id"` // ใช้ pointer เพราะอาจเป็น NULL
	Quantity   *int   `json:"quantity"`   // ใช้ pointer เพราะอาจเป็น NULL
	ReturnDate string `json:"return_date"`
}

type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

type Member struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	MLevel    string `json:"m_level"`
	TitleName string `json:"title_name"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
}

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
	RoleID   *int   `json:"role_id"` // ใช้ pointer เพราะอาจเป็น NULL
}

var db *sql.DB

func main() {
	// เชื่อมต่อฐานข้อมูล
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myapp")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// สร้าง Router
	router := mux.NewRouter()

	// Endpoint สำหรับแต่ละตาราง
	router.HandleFunc("/api/products", getProducts).Methods("GET")
	router.HandleFunc("/api/product_changes", getProductChanges).Methods("GET")
	router.HandleFunc("/api/reservations", getReservations).Methods("GET")
	router.HandleFunc("/api/returns", getReturns).Methods("GET")
	router.HandleFunc("/api/roles", getRoles).Methods("GET")
	router.HandleFunc("/api/members", getMembers).Methods("GET")
	router.HandleFunc("/api/users", getUsers).Methods("GET")

	// เริ่มเซิร์ฟเวอร์
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// ฟังก์ชันดึงข้อมูลจากแต่ละตาราง
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT product_id, product_name, category, quantity, barcode, stock_status, image_path, created_at, updated_at FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ProductID, &p.ProductName, &p.Category, &p.Quantity, &p.Barcode, &p.StockStatus, &p.ImagePath, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

func getProductChanges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT changes_id, change_type, old_quantity, new_quantity, changed_at, changed_by, created_at, product_id FROM product_changes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var changes []ProductChange
	for rows.Next() {
		var c ProductChange
		err := rows.Scan(&c.ChangesID, &c.ChangeType, &c.OldQuantity, &c.NewQuantity, &c.ChangedAt, &c.ChangedBy, &c.CreatedAt, &c.ProductID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		changes = append(changes, c)
	}
	json.NewEncoder(w).Encode(changes)
}

func getReservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT reserve_id, user_id, product_id, quantity, status, expires_at, created_at, actual_quantity, returned_quantity, updated_at FROM reservations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		err := rows.Scan(&r.ReserveID, &r.UserID, &r.ProductID, &r.Quantity, &r.Status, &r.ExpiresAt, &r.CreatedAt, &r.ActualQuantity, &r.ReturnedQuantity, &r.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reservations = append(reservations, r)
	}
	json.NewEncoder(w).Encode(reservations)
}

func getReturns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT return_id, reserve_id, user_id, product_id, quantity, return_date FROM returns")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var returns []Return
	for rows.Next() {
		var r Return
		err := rows.Scan(&r.ReturnID, &r.ReserveID, &r.UserID, &r.ProductID, &r.Quantity, &r.ReturnDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returns = append(returns, r)
	}
	json.NewEncoder(w).Encode(returns)
}

func getRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT role_id, role_name FROM roles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		err := rows.Scan(&r.RoleID, &r.RoleName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, r)
	}
	json.NewEncoder(w).Encode(roles)
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, username, password, m_level, title_name, name, surname FROM tbl_member")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		err := rows.Scan(&m.ID, &m.Username, &m.Password, &m.MLevel, &m.TitleName, &m.Name, &m.Surname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		members = append(members, m)
	}
	json.NewEncoder(w).Encode(members)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT user_id, username, email, phone, password, role, role_id FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.UserID, &u.Username, &u.Email, &u.Phone, &u.Password, &u.Role, &u.RoleID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}
