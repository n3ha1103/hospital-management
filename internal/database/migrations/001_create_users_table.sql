CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('receptionist', 'doctor')),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Insert default users
INSERT INTO users (username, email, password, role, first_name, last_name) VALUES
('admin_receptionist', 'receptionist@hospital.com', '$2a$10$YourHashedPasswordHere', 'receptionist', 'Admin', 'Receptionist'),
('dr_smith', 'dr.smith@hospital.com', '$2a$10$YourHashedPasswordHere', 'doctor', 'John', 'Smith');