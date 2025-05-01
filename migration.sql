BEGIN TRANSACTION;
DROP TABLE IF EXISTS "template_workouts";
CREATE TABLE IF NOT EXISTS "template_workouts" (
	"id"	INTEGER NOT NULL UNIQUE,
	"name"	TEXT,
	"user_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "real_sessions";
CREATE TABLE IF NOT EXISTS "real_sessions" (
	"id"	INTEGER NOT NULL UNIQUE,
	"template_session_id"	INTEGER,
	"real_workout_id"	INTEGER,
	"start_date"	TEXT,
	"finish_date"	TEXT,
	"user_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "template_exercises";
CREATE TABLE IF NOT EXISTS "template_exercises" (
	"id"	INTEGER NOT NULL UNIQUE,
	"name"	TEXT,
	"user_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "template_sessions";
CREATE TABLE IF NOT EXISTS "template_sessions" (
	"id"	INTEGER NOT NULL UNIQUE,
	"name"	TEXT,
	"user_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "template_sessions_template_exercises";
CREATE TABLE IF NOT EXISTS "template_sessions_template_exercises" (
	"template_session_id"	INTEGER,
	"template_exercise_id"	INTEGER
);
DROP TABLE IF EXISTS "template_workouts_template_sessions";
CREATE TABLE IF NOT EXISTS "template_workouts_template_sessions" (
	"template_workout_id"	INTEGER,
	"template_session_id"	INTEGER
);
DROP TABLE IF EXISTS "real_workouts";
CREATE TABLE IF NOT EXISTS "real_workouts" (
	"id"	INTEGER NOT NULL UNIQUE,
	"start_date"	TEXT,
	"finish_date"	TEXT,
	"template_workout_id"	INTEGER,
	"user_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "real_sets";
CREATE TABLE IF NOT EXISTS "real_sets" (
	"id"	INTEGER NOT NULL UNIQUE,
	"reps"	INTEGER,
	"weight"	INTEGER,
	"real_exercise_id"	INTEGER,
	"user_id"	INTEGER,
	"start_date"	TEXT,
	"finish_date"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "real_exercises";
CREATE TABLE IF NOT EXISTS "real_exercises" (
	"id"	INTEGER NOT NULL UNIQUE,
	"template_exercise_id"	INTEGER,
	"real_session_id"	INTEGER,
	"user_id"	INTEGER,
	"start_date"	TEXT,
	"finish_date"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "users";
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER NOT NULL UNIQUE,
	"username"	TEXT UNIQUE,
	"password"	TEXT,
	"session_id"	TEXT,
	"active_workout_id"	INTEGER,
	"salt"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
INSERT INTO "template_workouts" VALUES (1,'adf',16);
INSERT INTO "template_workouts" VALUES (2,'wwrwerwrwerrwrw',16);
INSERT INTO "template_workouts" VALUES (3,'push pull legs',16);
INSERT INTO "template_workouts" VALUES (6,'fadsfasfa',17);
INSERT INTO "template_workouts" VALUES (7,'rrrr',17);
INSERT INTO "template_workouts" VALUES (8,'srsrsrsrssrsr',17);
INSERT INTO "real_sessions" VALUES (157,8,34,'2025-04-27T00:07:24+02:00','2025-04-27T00:07:40+02:00',16);
INSERT INTO "real_sessions" VALUES (158,9,34,'2025-04-27T00:09:53+02:00','2025-04-27T00:20:03+02:00',16);
INSERT INTO "real_sessions" VALUES (159,12,36,NULL,NULL,17);
INSERT INTO "real_sessions" VALUES (160,13,38,NULL,NULL,17);
INSERT INTO "template_exercises" VALUES (1,'barbel',NULL);
INSERT INTO "template_exercises" VALUES (4,'aafff',16);
INSERT INTO "template_exercises" VALUES (5,'yeah',16);
INSERT INTO "template_exercises" VALUES (6,'dumbell curloo',16);
INSERT INTO "template_exercises" VALUES (7,'develope couche',16);
INSERT INTO "template_exercises" VALUES (8,'fasdfas',16);
INSERT INTO "template_exercises" VALUES (9,'adfsa',17);
INSERT INTO "template_sessions" VALUES (7,'mandddd',16);
INSERT INTO "template_sessions" VALUES (8,'2d',16);
INSERT INTO "template_sessions" VALUES (9,'3',16);
INSERT INTO "template_sessions" VALUES (10,'push',16);
INSERT INTO "template_sessions" VALUES (11,'pull',16);
INSERT INTO "template_sessions" VALUES (12,'asdfasf',17);
INSERT INTO "template_sessions" VALUES (13,'srsrsrsrsrs',17);
INSERT INTO "template_sessions_template_exercises" VALUES (9,4);
INSERT INTO "template_sessions_template_exercises" VALUES (10,7);
INSERT INTO "template_sessions_template_exercises" VALUES (11,6);
INSERT INTO "template_sessions_template_exercises" VALUES (11,5);
INSERT INTO "template_sessions_template_exercises" VALUES (7,4);
INSERT INTO "template_sessions_template_exercises" VALUES (7,5);
INSERT INTO "template_sessions_template_exercises" VALUES (8,4);
INSERT INTO "template_sessions_template_exercises" VALUES (8,5);
INSERT INTO "template_sessions_template_exercises" VALUES (12,9);
INSERT INTO "template_workouts_template_sessions" VALUES (2,7);
INSERT INTO "template_workouts_template_sessions" VALUES (2,8);
INSERT INTO "template_workouts_template_sessions" VALUES (3,10);
INSERT INTO "template_workouts_template_sessions" VALUES (3,11);
INSERT INTO "template_workouts_template_sessions" VALUES (1,8);
INSERT INTO "template_workouts_template_sessions" VALUES (1,9);
INSERT INTO "template_workouts_template_sessions" VALUES (6,12);
INSERT INTO "template_workouts_template_sessions" VALUES (8,13);
INSERT INTO "real_workouts" VALUES (34,'2025-04-27T00:00:49+02:00',NULL,1,16);
INSERT INTO "real_workouts" VALUES (35,'2025-05-02T00:10:04+02:00',NULL,5,17);
INSERT INTO "real_workouts" VALUES (36,'2025-05-02T00:13:18+02:00','2025-05-02T00:13:37+02:00',6,17);
INSERT INTO "real_workouts" VALUES (37,'2025-05-02T00:13:38+02:00','2025-05-02T00:15:05+02:00',7,17);
INSERT INTO "real_workouts" VALUES (38,'2025-05-02T00:15:05+02:00',NULL,8,17);
INSERT INTO "real_sets" VALUES (1,3,3,106,16,'1745605762',NULL);
INSERT INTO "real_sets" VALUES (2,4,3,106,16,'1745605774',NULL);
INSERT INTO "real_sets" VALUES (3,3,3,120,16,'2025-04-25T20:56:54+02:00',NULL);
INSERT INTO "real_sets" VALUES (4,34,3,120,16,'2025-04-25T20:56:56+02:00',NULL);
INSERT INTO "real_sets" VALUES (5,3,3,121,16,'2025-04-25T21:03:47+02:00',NULL);
INSERT INTO "real_sets" VALUES (6,3,4,130,16,'2025-04-26T16:30:57+02:00',NULL);
INSERT INTO "real_sets" VALUES (7,3,6,130,16,'2025-04-26T16:30:59+02:00',NULL);
INSERT INTO "real_sets" VALUES (8,3,3,135,16,'2025-04-26T18:16:27+02:00',NULL);
INSERT INTO "real_sets" VALUES (9,34,3,135,16,'2025-04-26T18:16:30+02:00',NULL);
INSERT INTO "real_sets" VALUES (10,3,3,135,16,'2025-04-26T18:17:39+02:00',NULL);
INSERT INTO "real_sets" VALUES (11,3,3,135,16,'2025-04-26T18:17:40+02:00',NULL);
INSERT INTO "real_sets" VALUES (12,3,3,135,16,'2025-04-26T18:17:41+02:00',NULL);
INSERT INTO "real_sets" VALUES (13,4,3,135,16,'2025-04-26T18:17:42+02:00',NULL);
INSERT INTO "real_sets" VALUES (14,3,55,135,16,'2025-04-26T23:52:30+02:00',NULL);
INSERT INTO "real_sets" VALUES (15,3,55,142,16,'2025-04-27T00:07:33+02:00',NULL);
INSERT INTO "real_sets" VALUES (16,5435,435345,144,16,'2025-04-27T00:19:59+02:00',NULL);
INSERT INTO "real_exercises" VALUES (114,4,139,16,NULL,'1745607227');
INSERT INTO "real_exercises" VALUES (115,5,139,16,NULL,'1745606505');
INSERT INTO "real_exercises" VALUES (116,4,140,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (117,4,141,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (118,5,141,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (119,4,142,16,NULL,'2025-04-25T20:55:55+02:00');
INSERT INTO "real_exercises" VALUES (120,5,142,16,NULL,'2025-04-25T20:56:56+02:00');
INSERT INTO "real_exercises" VALUES (121,4,143,16,NULL,'2025-04-25T21:03:48+02:00');
INSERT INTO "real_exercises" VALUES (122,5,143,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (123,4,144,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (124,4,145,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (125,5,145,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (126,4,146,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (127,4,147,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (128,5,147,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (129,4,148,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (130,7,149,16,NULL,'2025-04-26T16:31:00+02:00');
INSERT INTO "real_exercises" VALUES (131,6,150,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (132,5,150,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (133,7,151,16,NULL,'2025-04-26T17:33:37+02:00');
INSERT INTO "real_exercises" VALUES (134,6,152,16,NULL,'2025-04-26T17:36:54+02:00');
INSERT INTO "real_exercises" VALUES (135,5,152,16,NULL,'2025-04-26T23:52:32+02:00');
INSERT INTO "real_exercises" VALUES (136,7,153,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (137,6,154,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (138,5,154,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (139,4,155,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (140,5,155,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (141,4,156,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (142,4,157,16,NULL,'2025-04-27T00:07:35+02:00');
INSERT INTO "real_exercises" VALUES (143,5,157,16,NULL,NULL);
INSERT INTO "real_exercises" VALUES (144,4,158,16,NULL,'2025-04-27T00:20:00+02:00');
INSERT INTO "real_exercises" VALUES (145,9,159,17,NULL,NULL);
INSERT INTO "users" VALUES (10,'testuser','testpassword',NULL,NULL,NULL);
INSERT INTO "users" VALUES (11,'test','test','0364c0d4f672f0c780e6c13e88023053',NULL,NULL);
INSERT INTO "users" VALUES (12,'adsf','asdf','ec7bdc8892806552119a57a1a7ec42af',NULL,NULL);
INSERT INTO "users" VALUES (13,'sdafsf','asdfsa','c9ff9919408033b9602cb4181b21a328',NULL,NULL);
INSERT INTO "users" VALUES (14,'fas','fas','4953fc39180be44f1550b9da9dbe8b3c',NULL,NULL);
INSERT INTO "users" VALUES (15,'as','as','b10c4a08c19cc51c7b60bd63094ba019',NULL,NULL);
INSERT INTO "users" VALUES (16,'fa','fa','10237f2b479ba61fa9bbcce6ef98b768',34,NULL);
INSERT INTO "users" VALUES (17,'fd','$2a$10$vFWfm/MkM/HBjzVWtwxpdeNz6VX481AtB/XdYpX0zmr4scW28JOlK','80553b6d52db6325720baeb2af1e83ed',38,NULL);
INSERT INTO "users" VALUES (18,'ds','$2a$10$n5.4Q3YQEUhEd3KcpSMqOOZvS/MHYyr0oszO/4vnN5Sc8Ngq/ZD4i','a88fdddd6902ea7d9a9868519cc84f4e',NULL,'9c4887af1daef7198862bb295b3f395f');
COMMIT;
