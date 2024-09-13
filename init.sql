CREATE TABLE employee (
    id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id VARCHAR(100) REFERENCES organization(id) ON DELETE CASCADE,
    user_id VARCHAR(100) REFERENCES employee(id) ON DELETE CASCADE
);


-------------------------------------
---DDL для таблиц тендеров и бидов---
-------------------------------------

CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');

CREATE TYPE tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');

CREATE TABLE tenders (
    id VARCHAR(100) PRIMARY KEY,  -- tenderId
    name VARCHAR(100) NOT NULL,  -- tenderName
    description VARCHAR(500) NOT NULL,  -- tenderDescription
    service_type tender_service_type NOT NULL,  -- tenderServiceType --- (ENUM)
    status tender_status NOT NULL,  -- tenderStatus --- (ENUM)
    organization_id VARCHAR(100) NOT NULL,  -- organizationId
    version INT NOT NULL DEFAULT 1 CHECK (version >= 1),  -- tenderVersion
    created_at TIMESTAMP NOT NULL  -- createdAt
);

CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');

CREATE TYPE bid_decision AS ENUM ('Approved', 'Rejected'); --- In openapi.yaml there was no info about bid_decision's usage, so I took responsibility upon myself to add it to bid_reviews

CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');

CREATE TABLE bids (
    id VARCHAR(100) PRIMARY KEY,  -- bidId
    name VARCHAR(100) NOT NULL,  -- bidName
    description VARCHAR(500) NOT NULL,  -- bidDescription
    status bid_status NOT NULL,  -- bidStatus --- (ENUM)
    tender_id VARCHAR(100) REFERENCES tenders(id),  -- tenderId (FK)
    author_type bid_author_type NOT NULL,  -- bidAuthorType --- (ENUM)
    author_id VARCHAR(100) NOT NULL REFERENCES employee(id),  -- bidAuthorId (FK)
    version INT NOT NULL DEFAULT 1,  -- bidVersion
    created_at TIMESTAMP NOT NULL  -- createdAt
);

CREATE TABLE bid_reviews (
    id VARCHAR(100) PRIMARY KEY,  -- bidReviewId
    description VARCHAR(1000) NOT NULL,  -- bidReviewDescription
    created_at TIMESTAMP NOT NULL,  -- createdAt
    bid_id VARCHAR(100) REFERENCES bids(id),  -- FK on bid's ID
    decision bid_decision NOT NULL -- bid_decision --- (ENUM)
);

------------------------
---DML with test data---
------------------------

INSERT INTO employee (username, first_name, last_name)
VALUES 
('john_doe', 'John', 'Doe'),
('jane_smith', 'Jane', 'Smith'),
('alex_brown', 'Alex', 'Brown'),
('lisa_jones', 'Lisa', 'Jones');

INSERT INTO organization (name, description, type)
VALUES 
('Tech Innovators', 'A technology company specializing in software solutions', 'LLC'),
('Green Builders', 'Eco-friendly construction company', 'IE'),
('Fast Delivery Co', 'Delivery services across the country', 'JSC');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES 
(1, 1),  -- John Doe is responsible for Tech Innovators
(2, 2),  -- Jane Smith is responsible for Green Builders
(3, 3);  -- Alex Brown is responsible for Fast Delivery Co

INSERT INTO tenders (id, name, description, service_type, status, organization_id, created_at)
VALUES 
('TND001', 'New Office Construction', 'Construction of a new office building', 'Construction', 'Created', '1', '2023-09-01 10:00:00'),
('TND002', 'Warehouse Delivery', 'Delivery of materials to the warehouse', 'Delivery', 'Published', '3', '2023-09-02 12:00:00'),
('TND003', 'Product Manufacture', 'Manufacture of electronic devices', 'Manufacture', 'Closed', '1', '2023-09-05 14:00:00');

INSERT INTO bids (id, name, description, status, tender_id, author_type, author_id, created_at)
VALUES 
('BID001', 'Bid for Office Construction', 'Our proposal for the office construction', 'Published', 'TND001', 'Organization', '1', '2023-09-03 09:00:00'),
('BID002', 'Delivery Bid', 'Proposal for the warehouse delivery', 'Created', 'TND002', 'User', '2', '2023-09-04 11:00:00'),
('BID003', 'Manufacturing Proposal', 'Proposal for product manufacturing', 'Approved', 'TND003', 'Organization', '3', '2023-09-06 16:00:00');

INSERT INTO bid_reviews (id, description, created_at, bid_id, decision)
VALUES 
('REV001', 'Approved the construction bid for meeting all requirements', '2023-09-07 10:00:00', 'BID001', 'Approved'),
('REV002', 'Rejected the delivery bid due to pricing issues', '2023-09-08 12:00:00', 'BID002', 'Rejected'),
('REV003', 'Approved the manufacturing proposal as it fits the timeline', '2023-09-09 14:00:00', 'BID003', 'Approved');
