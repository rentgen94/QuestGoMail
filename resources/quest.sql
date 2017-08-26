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
DROP TABLE IF EXISTS LabyrinthRoomLink;
DROP TABLE IF EXISTS LabyrinthActionLink;
DROP TABLE IF EXISTS ActionInteractiveLink;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id       SERIAL PRIMARY KEY,
  name     VARCHAR(256) UNIQUE NOT NULL,
  password VARCHAR(256)        NOT NULL
);

CREATE TABLE Room (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  description VARCHAR(500) NOT NULL
);

CREATE TABLE Door (
  id           SERIAL PRIMARY KEY,
  room1        INT REFERENCES Room (id) NOT NULL,
  room2        INT REFERENCES Room (id) NOT NULL,
  name         VARCHAR(100) UNIQUE      NOT NULL,
  isAccessible BOOLEAN                  NOT NULL,
  UNIQUE (room1, room2)
);
CREATE INDEX door_room1_idx
  ON Door (room1);
CREATE INDEX door_room2_idx
  ON Door (room2);

CREATE TABLE Item (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) UNIQUE NOT NULL,
  description VARCHAR(500)        NOT NULL,
  size        INT                 NOT NULL
);

CREATE TABLE Slot (
  id           SERIAL PRIMARY KEY,
  room         INT REFERENCES Room (id),
  name         VARCHAR(100),
  capacity     INT,
  isAccessible BOOLEAN
);
CREATE INDEX slot_room_idx
  ON Slot (room);

CREATE TABLE Interactive (
  id           SERIAL PRIMARY KEY,
  room         INT REFERENCES Room (id) NOT NULL,
  name         VARCHAR(100)             NOT NULL,
  description  VARCHAR(500)             NOT NULL,
  isAccessible BOOLEAN                  NOT NULL,
  args         VARCHAR(100)             NOT NULL --this is an argument string. Arguments has to be divided with ' '
);
CREATE INDEX interactive_room_idx
  ON Interactive (room);

CREATE TABLE Action (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(100)  NOT NULL,
  resultCode INT DEFAULT 0 NOT NULL,
  resultMsg  VARCHAR(500)  NOT NULL
);

CREATE TABLE Labyrinth (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) UNIQUE      NOT NULL,
  description VARCHAR(500)             NOT NULL,
  startRoom   INT REFERENCES Room (id) NOT NULL
);
CREATE INDEX labyrinth_room_idx
  ON Labyrinth (startRoom);

CREATE TABLE InteractiveObjectNeed (
  id          SERIAL PRIMARY KEY,
  interactive INT REFERENCES Interactive (id) NOT NULL,
  item        INT REFERENCES Item (id)        NOT NULL,
  UNIQUE (interactive, item)
);
CREATE INDEX interactiveObjectNeed_idx
  ON InteractiveObjectNeed (interactive, item);

CREATE TABLE ActionInteractiveLink (
  id          SERIAL PRIMARY KEY,
  action      INT REFERENCES Action (id)      NOT NULL,
  interactive INT REFERENCES Interactive (id) NOT NULL,
  UNIQUE (action, interactive)
);
CREATE INDEX actionInteractiveLink_idx
  ON ActionInteractiveLink (action, interactive);

CREATE TABLE ActionInteractiveSwitch (
  id          SERIAL PRIMARY KEY,
  action      INT REFERENCES Action (id)      NOT NULL,
  interactive INT REFERENCES Interactive (id) NOT NULL,
  newState    BOOLEAN                         NOT NULL,
  UNIQUE (action, interactive)
);
CREATE INDEX actionInteractiveSwitch_idx
  ON ActionInteractiveSwitch (action, interactive);

CREATE TABLE ActionDoorSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id) NOT NULL,
  door     INT REFERENCES Door (id)   NOT NULL,
  newState BOOLEAN,
  UNIQUE (action, door)
);
CREATE INDEX actionDoorSwitch_idx
  ON ActionDoorSwitch (action, door);

CREATE TABLE ActionSlotSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id) NOT NULL,
  slot     INT REFERENCES Slot (id)   NOT NULL,
  newState BOOLEAN,
  UNIQUE (action, slot)
);
CREATE INDEX actionSlotSwitch_idx
  ON ActionSlotSwitch (action, slot);

CREATE TABLE SlotItemLink (
  id   SERIAL PRIMARY KEY,
  slot INT REFERENCES Slot (id) NOT NULL,
  item INT REFERENCES Item (id) NOT NULL,
  UNIQUE (slot, item)
);
CREATE INDEX SlotItemLink_idx
  ON SlotItemLink (slot, item);

CREATE TABLE LabyrinthRoomLink (
  id        SERIAL PRIMARY KEY,
  room      INT REFERENCES Room (id)      NOT NULL,
  labyrinth INT REFERENCES Labyrinth (id) NOT NULL,
  UNIQUE (room, labyrinth)
);
CREATE INDEX labyrinthRoomLink_idx
  ON LabyrinthRoomLink (room, labyrinth);

CREATE TABLE LabyrinthActionLink (
  id        SERIAL PRIMARY KEY,
  action    INT REFERENCES Action (id)    NOT NULL,
  labyrinth INT REFERENCES Labyrinth (id) NOT NULL,
  UNIQUE (action, labyrinth)
);
CREATE INDEX labyrinthActionLink_idx
  ON LabyrinthActionLink (action, labyrinth);
