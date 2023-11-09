-- COTDB Startup Script
CREATE DATABASE cotdb;
USE cotdb;

-- Table creation
CREATE TABLE Airports
(IATA_Code varchar(4) NOT NULL,
 comment varchar(255),
 CONSTRAINT airports_pk PRIMARY KEY (IATA_Code));

CREATE TABLE Roles
(role_name varchar(10) NOT NULL,
 description text,
 CONSTRAINT roles_pk PRIMARY KEY (role_name));

CREATE TABLE Members
(phone_number varchar(11) NOT NULL,
 name varchar(30),
 added_by varchar(30),
 notify int,
 isValid int,
 timestamp bigint,
 CONSTRAINT members_pk PRIMARY KEY (phone_number));

CREATE TABLE Flies
(phone_number varchar(11) NOT NULL,
 IATA_Code varchar(4) NOT NULL,
 CONSTRAINT flies_fk_members FOREIGN KEY (phone_number) REFERENCES Members (phone_number),
 CONSTRAINT flies_fk_airports FOREIGN KEY (IATA_Code) REFERENCES Airports (IATA_Code));

CREATE TABLE Wants
(phone_number varchar(11) NOT NULL,
 role_name varchar(10) NOT NULL,
 CONSTRAINT wants_fk_members FOREIGN KEY (phone_number) REFERENCES Members (phone_number),
 CONSTRAINT wants_fk_roles FOREIGN KEY (role_name) REFERENCES Roles (role_name));

 CREATE TABLE Requester
 (timestamp bigint NOT NULL,
  phone_number varchar(11) NOT NULL,
  request_message text NOT NULL,
  CONSTRAINT requester_pk PRIMARY KEY (timestamp),
  CONSTRAINT requester_fk_members FOREIGN KEY (phone_number) REFERENCES Members (phone_number));

CREATE TABLE Contacts
(requester_phone varchar(11) NOT NULL,
 contacted_phone varchar(11) NOT NULL,
 timestamp bigint NOT NULL,
 CONSTRAINT contacts_pk PRIMARY KEY (contacted_phone));

CREATE TABLE Deferred (
    requester_phone varchar(11) NOT NULL,
    contacted_phone varchar(11) NOT NULL,
    request_message text NOT NULL,
    timestamp bigint NOT NULL,
    CONSTRAINT deferred_pk PRIMARY KEY (requester_phone, contacted_phone, timestamp)
);

-- Populate static tables
INSERT INTO Roles (role_name, description) VALUES ("ADIS", "Aerial Digital Imaging System Operator");
INSERT INTO Roles (role_name, description) VALUES ("AOBD", "Air Operations Branch Director");
INSERT INTO Roles (role_name, description) VALUES ("AP", "Aerial Photographer");
INSERT INTO Roles (role_name, description) VALUES ("CD", "Counterdrug");
INSERT INTO Roles (role_name, description) VALUES ("CERT", "Community Emergency Response Team");
INSERT INTO Roles (role_name, description) VALUES ("CISM", "Critical Incident Stress Management Personnel");
INSERT INTO Roles (role_name, description) VALUES ("CSSCS", "Chaplain Support Specialist (CAP Support Rating)");
INSERT INTO Roles (role_name, description) VALUES ("CUL", "Communications Unit Leader");
INSERT INTO Roles (role_name, description) VALUES ("DAARTO", "Domestic Operations Awareness and Assessment Response Tool Operator");
INSERT INTO Roles (role_name, description) VALUES ("DAARTU", "Domestic Operations Awareness and Assessment Repsonse Tool User");
INSERT INTO Roles (role_name, description) VALUES ("FASC", "Finance/Admin Section Chief");
INSERT INTO Roles (role_name, description) VALUES ("FLM", "Flight Line Marshaller");
INSERT INTO Roles (role_name, description) VALUES ("FLS", "Flight Line Supervisor");
INSERT INTO Roles (role_name, description) VALUES ("FRO", "Flight Release Officer");
INSERT INTO Roles (role_name, description) VALUES ("GDB", "Ground Branch Director");
INSERT INTO Roles (role_name, description) VALUES ("GES", "General Emergency Services");
INSERT INTO Roles (role_name, description) VALUES ("GFMC", "Surrogate Unmanned Aerial System Green Flag Mission Coordinator");
INSERT INTO Roles (role_name, description) VALUES ("GFMP", "Surrogate Unmanned Aerial System Green Flag Mission Pilot");
INSERT INTO Roles (role_name, description) VALUES ("GFSO", "Surrogate Unmanned Aerial System Green Flag Sensor Operator");
INSERT INTO Roles (role_name, description) VALUES ("GTL", "Ground Team Leader");
INSERT INTO Roles (role_name, description) VALUES ("GTM1", "Ground Team Member Level 1");
INSERT INTO Roles (role_name, description) VALUES ("GTM2", "Ground Team Member Level 2");
INSERT INTO Roles (role_name, description) VALUES ("GTM3", "Ground Team Member Level 3");
INSERT INTO Roles (role_name, description) VALUES ("IC1", "Incident Commander Level 1");
INSERT INTO Roles (role_name, description) VALUES ("IC2", "Incident Commander Level 2");
INSERT INTO Roles (role_name, description) VALUES ("IC3", "Incident Commander Level 3");
INSERT INTO Roles (role_name, description) VALUES ("ICE", "Incident Commander");
INSERT INTO Roles (role_name, description) VALUES ("LO", "Liaison Officer");
INSERT INTO Roles (role_name, description) VALUES ("LSC", "Logistics Sectio Chief");
INSERT INTO Roles (role_name, description) VALUES ("MCCS", "Mission Chaplian (CAP Support Rating)");
INSERT INTO Roles (role_name, description) VALUES ("MCDS", "Mission Chaplain (Disaster Support Rating)");
INSERT INTO Roles (role_name, description) VALUES ("MFC", "Mountain Flying Certification");
INSERT INTO Roles (role_name, description) VALUES ("MO", "Mission Observer");
INSERT INTO Roles (role_name, description) VALUES ("MP", "Mission Pilot");
INSERT INTO Roles (role_name, description) VALUES ("MRO", "Mission Radio Operator");
INSERT INTO Roles (role_name, description) VALUES ("MS", "Mission Scanner");
INSERT INTO Roles (role_name, description) VALUES ("MSA", "Mission Staff Assistant");
INSERT INTO Roles (role_name, description) VALUES ("MSO", "Mission Safety Officer");
INSERT INTO Roles (role_name, description) VALUES ("OSC", "Operations Section Chief");
INSERT INTO Roles (role_name, description) VALUES ("PAX", "Passenger");
INSERT INTO Roles (role_name, description) VALUES ("PIO", "Public Information Officer");
INSERT INTO Roles (role_name, description) VALUES ("PODC", "Point of Distribution Course");
INSERT INTO Roles (role_name, description) VALUES ("PSC", "Planning Section Chief");
INSERT INTO Roles (role_name, description) VALUES ("SFGC", "Shelter Field Guide Course");
INSERT INTO Roles (role_name, description) VALUES ("SFRO", "Senior Flight Release Officer");
INSERT INTO Roles (role_name, description) VALUES ("TMP", "Transport Mission Pilot");
INSERT INTO Roles (role_name, description) VALUES ("UAO", "Unit Alert Officer");
INSERT INTO Roles (role_name, description) VALUES ("UASMP", "sUAS Mission Pilot");
INSERT INTO Roles (role_name, description) VALUES ("UASP", "sUAS Pilot");
INSERT INTO Roles (role_name, description) VALUES ("UAST", "sUAS Technician");
INSERT INTO Roles (role_name, description) VALUES ("UDF", "Urban Direction Finding Team");

INSERT INTO Airports (IATA_Code, comment) VALUES ("KALW", "Walla Walla Regional Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KAWO", "Arlington Municipal Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KBFI", "King Couny International Airport-Boeing Field");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KBLI", "Bellingham International Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KBVS", "Skagit Regional Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KCLM", "William R Fairchild International Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KCLS", "Chehalis-Centralia Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KDEW", "Deer Park Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KEAT", "Pangborn Memorial Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KELN", "Bowers Field Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KEPH", "Ephrata Municipal Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KGEG", "Spokane International Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KMWH", "Grant County International Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KOLM", "Olympia Regional Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KOMK", "Omak Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KPAE", "Snohomish County Airport (Paine Field)");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KPLU", "Pierce County Airport - Thun Field");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KPSC", "Tri-Cities Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KPUW", "Pullman-Moscow Regional Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KPWT", "Bremerton National Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KRNT", "Renton Municipal Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KS50", "Auburn Municipal Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KTIW", "Tacoma Narrows Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KVUO", "Pearson Field Airport");
INSERT INTO Airports (IATA_Code, comment) VALUES ("KYKM", "Yakima Air Terminal-McAllister Field");

-- Add a couple non-existant test members
INSERT INTO Members (phone_number, name, added_by, notify, isValid, timestamp) VALUES ("14564564567", "Moe Kibley", "10987654321", 1, 1, 0);
INSERT INTO Members (phone_number, name, added_by, notify, isValid, timestamp) VALUES ("12342342345", "Ape Kibley", "10987654321", 1, 1, 0);
INSERT INTO Members (phone_number, name, added_by, notify, isValid, timestamp) VALUES ("11234567890", "Default User", "10987654321", 1, 1, 0);

INSERT INTO Flies (phone_number, IATA_Code) VALUES ("14564564567", "KBLI");
INSERT INTO Flies (phone_number, IATA_Code) VALUES ("12342342345", "KBLI");

INSERT INTO Wants (phone_number, role_name) VALUES ("14564564567", "MO");
INSERT INTO Wants (phone_number, role_name) VALUES ("12342342345", "AP");
INSERT INTO Wants (phone_number, role_name) VALUES ("12342342345", "MO");
