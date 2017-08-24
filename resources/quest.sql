DROP TABLE IF EXISTS Room CASCADE;
DROP TABLE IF EXISTS Door CASCADE;
DROP TABLE IF EXISTS Item CASCADE;
DROP TABLE IF EXISTS Interactive CASCADE;
DROP TABLE IF EXISTS Action CASCADE;
DROP TABLE IF EXISTS Slot CASCADE;
DROP TABLE IF EXISTS Labyrinth CASCADE;
DROP TABLE IF EXISTS InteractiveObjectNeed;
DROP TABLE IF EXISTS ActionInteractiveSwitch;
DROP TABLE IF EXISTS ActionDoorSwitch;
DROP TABLE IF EXISTS ActionSlotSwitch;
DROP TABLE IF EXISTS SlotItemLink;
DROP TABLE IF EXISTS RoomScenarioLink;

CREATE TABLE Room (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100),
  description VARCHAR(500)
);

CREATE TABLE Door (
  id           SERIAL PRIMARY KEY,
  room1        INT REFERENCES Room (id),
  room2        INT REFERENCES Room (id),
  name         VARCHAR(100) UNIQUE,
  isAccessible BOOLEAN,
  UNIQUE (room1, room2)
);

CREATE TABLE Item (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) UNIQUE,
  description VARCHAR(500),
  size        INT
);

CREATE TABLE Slot (
  id           SERIAL PRIMARY KEY,
  room         INT REFERENCES Room (id),
  name         VARCHAR(100),
  capacity     INT,
  isAccessible BOOLEAN
);

CREATE TABLE Interactive (
  id           SERIAL PRIMARY KEY,
  room         INT REFERENCES Room (id),
  name         VARCHAR(100),
  description  VARCHAR(500),
  isAccessible BOOLEAN,
  args         VARCHAR(25) []
);

CREATE TABLE Action (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(100),
  resultCode INT DEFAULT 0,
  resultMsg  VARCHAR(500)
);

CREATE TABLE Labyrinth (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE
);

CREATE TABLE InteractiveObjectNeed (
  id          SERIAL PRIMARY KEY,
  interactive INT REFERENCES Interactive (id),
  item        INT REFERENCES Item (id)
);

CREATE TABLE ActionInteractiveSwitch (
  id          SERIAL PRIMARY KEY,
  action      INT REFERENCES Action (id),
  interactive INT REFERENCES Interactive (id),
  newState    BOOLEAN
);

CREATE TABLE ActionDoorSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id),
  door     INT REFERENCES Door (id),
  newState BOOLEAN
);

CREATE TABLE ActionSlotSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id),
  slot     INT REFERENCES Slot (id),
  newState BOOLEAN
);

CREATE TABLE SlotItemLink (
  id   SERIAL PRIMARY KEY,
  slot INT REFERENCES Slot (id),
  item INT REFERENCES Item (id)
);
