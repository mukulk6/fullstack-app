# 🛒 Product Management Dashboard — React + Go + PostgreSQL

This is a full-stack web application built using **React** (frontend), **GoLang (Gin)** (backend), and **PostgreSQL** (database). The app supports **role-based authentication** for Users and Admins and allows for full CRUD operations on products by Admins, while Users can only view product details.

---

## 🚀 Features

### 👤 User Roles

- **User**
  - Sign up and log in
  - View list of available products
  - View product details on click

- **Admin**
  - Full access to product management:
    - ✅ Create
    - ✏️ Edit
    - ❌ Delete
  - View list of products created by each Admin user

---

## 🛠️ Tech Stack

| Layer      | Tech Used           |
|------------|---------------------|
| Frontend   | React.js + Material UI |
| Backend    | GoLang (Gin Framework) |
| Database   | PostgreSQL          |
| Auth       | JWT-based authentication |
| Styling    | Material UI (MUI)   |

---

## 📦 Setup Instructions

### Prerequisites

- Node.js & npm
- Go (1.18+)
- PostgreSQL

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name
