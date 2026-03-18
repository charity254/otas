# Otas – Full Product & Development Guide (Refined)

**Production-level fintech savings application for the Kenyan market**
Backend: Go | Frontend: React

---

# 1. Product Overview

## 1.1 Vision

Otas helps users **save money automatically** by rounding off everyday transactions and allocating the difference toward savings goals.

---

## 1.2 Target Users

* Students managing small daily expenses
* Young professionals with inconsistent saving habits
* Mobile-money-heavy users (e.g. M-Pesa users)

---

## 1.3 Problem Statement

Many users struggle to save consistently because:

* Savings require discipline
* Small amounts are often ignored
* No automation exists for micro-savings

---

## 1.4 Solution

Otas automates savings by:

* Rounding transactions
* Saving the difference
* Allocating savings toward goals

---

# 2. User Journey (Core Flow)

1. User enters phone number
2. Receives OTP
3. Logs in successfully
4. Views dashboard
5. Enters transaction amount
6. System calculates round-up
7. Round-up stored as savings
8. Background worker allocates to goals
9. UI updates in real-time

---

# 3. Product Requirements

## 3.1 Core Features

The system MUST support:

* OTP-based authentication
* Wallet per user
* Manual transaction entry
* Automatic round-up savings
* Savings goals with tracking
* Background processing
* Real-time notifications (WebSocket)
* Financial dashboard

---

## 3.2 Functional Requirements

### Authentication

* Login via phone number
* OTP expires within 5 minutes
* JWT issued after verification

---

### Transactions

* User submits amount
* System calculates round-up
* Round-up stored separately

---

### Savings

* Users create goals
* Round-ups allocated automatically
* Processing must be asynchronous

---

### Financial Integrity

* All money flows through ledger
* No direct balance updates allowed

---

### Frontend Requirements

Must display:

* Wallet balance
* Total savings
* Goals progress
* Transaction history

Must support real-time updates

---

## 3.3 Non-Functional Requirements

* RESTful API
* Concurrency-safe system
* Full input validation
* JWT-secured routes
* Clean architecture

---

# 4. MVP Definition (Phase 1 Scope)

To avoid overbuilding, MVP includes ONLY:

* OTP authentication
* Transaction input
* Round-up calculation
* Savings display (basic)
* Wallet balance

Everything else (goals, WS, worker) comes later.

---

# 5. API Design (Initial)

```
POST /auth/request-otp
POST /auth/verify-otp

GET  /wallet

POST /transactions
GET  /transactions

POST /goals
GET  /goals
```

---

# 6. Database Design (Simplified)

```
users(id, phone, created_at)

wallets(id, user_id, balance)

transactions(id, user_id, amount, created_at)

roundups(id, transaction_id, amount, processed)

goals(id, user_id, target_amount, current_amount)

ledger(id, user_id, amount, type, reference_id, created_at)
```

---

# 7. System Architecture Overview

### Backend (Go)

* REST API
* Business logic
* Background worker
* WebSocket server

### Frontend (React)

* UI + state management
* API integration
* Real-time updates

### Database

* PostgreSQL

---

# 8. Deliverables Overview

| #  | Deliverable          | Owner    | Phase |
| -- | -------------------- | -------- | ----- |
| 1  | Project setup        | Backend  | 0     |
| 2  | DB schema            | Backend  | 1     |
| 3  | Auth system          | Backend  | 1     |
| 4  | Transactions logic   | Backend  | 1     |
| 5  | Goals API            | Backend  | 2     |
| 6  | Worker               | Backend  | 2     |
| 7  | Ledger system        | Backend  | 2     |
| 8  | WebSocket            | Backend  | 2     |
| 9  | Frontend setup       | Frontend | 2     |
| 10 | Dashboard UI         | Frontend | 2     |
| 11 | API integration      | Frontend | 2     |
| 12 | Validation & logging | Backend  | 3     |
| 13 | Testing              | Backend  | 3     |
| 14 | Deployment           | All      | 4     |

---

# 9. Phase Breakdown

## Phase 0 — Setup

* Project structure
* Config management
* DB connection
* Server bootstrap

---

## Phase 1 — Core Backend (MVP)

* Auth (OTP + JWT)
* Transactions
* Round-up logic

---

## Phase 2 — Savings Engine + Frontend

* Goals system
* Background worker
* Ledger enforcement
* WebSocket updates
* Full frontend UI

---

## Phase 3 — Hardening

* Validation
* Rate limiting
* Logging
* Error handling
* Testing

---

## Phase 4 — Production

* Deployment
* Monitoring
* Stability checks

---

# 10. Critical System Rules

### 10.1 Round-Up Formula

```
roundup = ceil(amount/100)*100 - amount
```

---

### 10.2 Ledger Rule (STRICT)

* Every balance change MUST go through ledger
* No shortcuts allowed

---

### 10.3 Idempotency

* Prevent double-processing of round-ups
* Use `processed` flag

---

# 11. Error Handling Strategy

System must handle:

* Invalid OTP
* Expired OTP
* Duplicate transactions
* Worker failures

Approach:

* Retry logic
* Clear error messages
* Logging for debugging

---

# 12. Security Considerations

* OTP expiration enforced
* Rate limiting on auth endpoints
* JWT expiration
* HTTPS required
* Input validation on all endpoints

---

# 13. Final Definition of Done

The system is complete when:

* Users log in via OTP
* Transactions are recorded
* Round-ups are accurate
* Savings update automatically
* Ledger reflects all activity
* UI updates in real-time
* System is deployed and stable

---

# 14. Future Enhancements (Post-MVP)

* M-Pesa integration
* Auto-transaction detection
* AI-based saving suggestions
* Multi-goal allocation strategies
* Analytics dashboard

---
