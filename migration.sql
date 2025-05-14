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
INSERT INTO "template_workouts" VALUES (10,'alpha push',21);
INSERT INTO "real_sessions" VALUES (162,15,40,NULL,NULL,21);
INSERT INTO "real_sessions" VALUES (163,17,40,NULL,NULL,21);
INSERT INTO "real_sessions" VALUES (164,15,41,'2025-05-03T17:12:25+02:00','2025-05-03T17:12:29+02:00',21);
INSERT INTO "real_sessions" VALUES (165,17,41,NULL,NULL,21);
INSERT INTO "real_sessions" VALUES (166,18,41,NULL,NULL,21);
INSERT INTO "template_exercises" VALUES (11,'developpe couche',21);
INSERT INTO "template_exercises" VALUES (12,'calin arbre',21);
INSERT INTO "template_exercises" VALUES (13,'developpe militaire',21);
INSERT INTO "template_exercises" VALUES (14,'rear delt machine',21);
INSERT INTO "template_exercises" VALUES (15,'triceps skullcrusher',21);
INSERT INTO "template_exercises" VALUES (16,'triceps poulie prise inversee',21);
INSERT INTO "template_exercises" VALUES (22,'df',21);
INSERT INTO "template_exercises" VALUES (23,'develoope couche incline',21);
INSERT INTO "template_exercises" VALUES (24,'dips/dips assis superset',21);
INSERT INTO "template_exercises" VALUES (25,'arnold press',21);
INSERT INTO "template_exercises" VALUES (26,'triceps poulie',21);
INSERT INTO "template_exercises" VALUES (27,'oiseau avec oiseau leger pour tuer',21);
INSERT INTO "template_exercises" VALUES (28,'exercices poulie rehab',21);
INSERT INTO "template_exercises" VALUES (29,'tirage vertical large',21);
INSERT INTO "template_exercises" VALUES (30,'soulever de terre',21);
INSERT INTO "template_exercises" VALUES (31,'biceps prise marteau en diago',21);
INSERT INTO "template_exercises" VALUES (32,'bucheron',21);
INSERT INTO "template_exercises" VALUES (33,'poulie biceps avec la barre droite',21);
INSERT INTO "template_exercises" VALUES (34,'avant bras',21);
INSERT INTO "template_exercises" VALUES (35,'shrugs',21);
INSERT INTO "template_exercises" VALUES (36,'cou en partant',21);
INSERT INTO "template_sessions" VALUES (15,'push 1',21);
INSERT INTO "template_sessions" VALUES (17,'push 2',21);
INSERT INTO "template_sessions" VALUES (18,'pull 1',21);
INSERT INTO "template_sessions_template_exercises" VALUES (15,11);
INSERT INTO "template_sessions_template_exercises" VALUES (15,12);
INSERT INTO "template_sessions_template_exercises" VALUES (15,13);
INSERT INTO "template_sessions_template_exercises" VALUES (15,14);
INSERT INTO "template_sessions_template_exercises" VALUES (15,15);
INSERT INTO "template_sessions_template_exercises" VALUES (15,16);
INSERT INTO "template_sessions_template_exercises" VALUES (16,22);
INSERT INTO "template_sessions_template_exercises" VALUES (17,23);
INSERT INTO "template_sessions_template_exercises" VALUES (17,24);
INSERT INTO "template_sessions_template_exercises" VALUES (17,25);
INSERT INTO "template_sessions_template_exercises" VALUES (17,26);
INSERT INTO "template_sessions_template_exercises" VALUES (17,27);
INSERT INTO "template_sessions_template_exercises" VALUES (17,28);
INSERT INTO "template_sessions_template_exercises" VALUES (18,29);
INSERT INTO "template_sessions_template_exercises" VALUES (18,30);
INSERT INTO "template_sessions_template_exercises" VALUES (18,31);
INSERT INTO "template_sessions_template_exercises" VALUES (18,32);
INSERT INTO "template_sessions_template_exercises" VALUES (18,33);
INSERT INTO "template_sessions_template_exercises" VALUES (18,34);
INSERT INTO "template_sessions_template_exercises" VALUES (18,34);
INSERT INTO "template_sessions_template_exercises" VALUES (18,36);
INSERT INTO "template_workouts_template_sessions" VALUES (10,15);
INSERT INTO "template_workouts_template_sessions" VALUES (10,17);
INSERT INTO "template_workouts_template_sessions" VALUES (10,18);
INSERT INTO "real_workouts" VALUES (40,'2025-05-03T16:39:25+02:00','2025-05-03T16:40:33+02:00',10,21);
INSERT INTO "real_workouts" VALUES (41,'2025-05-03T16:40:33+02:00',NULL,10,21);
INSERT INTO "real_exercises" VALUES (147,11,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (148,12,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (149,13,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (150,14,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (151,15,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (152,16,162,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (153,23,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (154,24,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (155,25,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (156,26,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (157,27,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (158,28,163,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (159,11,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (160,12,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (161,13,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (162,14,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (163,15,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (164,16,164,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (165,23,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (166,24,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (167,25,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (168,26,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (169,27,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (170,28,165,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (171,29,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (172,30,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (173,31,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (174,32,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (175,33,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (176,34,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (177,34,166,21,NULL,NULL);
INSERT INTO "real_exercises" VALUES (178,36,166,21,NULL,NULL);
INSERT INTO "users" VALUES (21,'fd','$2a$10$4G5XDI7uTRdHzFbvteDEQupUnKBE8f257w3Cods7nnrahHtQWqqAW','4dfb4f68a428ee19e6d4c93888c90153',41,'f87f1f1bc5269ce54f691db18cf0f066');
COMMIT;
