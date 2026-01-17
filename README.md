# RabbitMQ Worker Service

Service สำหรับรับ Event (Socket Events) จาก RabbitMQ แล้วบันทึกข้อมูลลง Azure PostgreSQL Database โดยทำงานเป็น Background Worker ที่เขียนด้วย Go

## Requirements
- Docker & Docker Compose

## Installation & Run

1. **Config:** แก้ไขไฟล์ `.env` เพื่อตั้งค่า Database และ RabbitMQ User/Pass
2. **Start:** รันคำสั่งด้านล่างเพื่อเริ่มระบบ
   ```bash
   docker-compose up -d --build
   ```
3. **Logs:** ดู Logs ของ Worker
   ```bash
   docker-compose logs -f worker
   ```

## Project Structure
```
rabbitMQ/
├── worker/
│   ├── database/    # DB Connection & Auto Migration
│   ├── handlers/    # Business Logic & SQL Queries
│   ├── models/      # Struct Definitions
│   ├── rabbitmq/    # Consumer Loop Connection
│   ├── utils/       # Logging & Helpers
│   └── main.go      # App Entry Point
└── docker-compose.yml
```

## Testing
- **RabbitMQ Management UI:** [http://localhost:15672](http://localhost:15672)
- **Queue Name:** `socket_events`
- **Default User/Pass:** `user` / `password` (หรือตามค่าใน .env)
