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
CREATE INDEX door_room1_idx
  ON Door (room1);
CREATE INDEX door_room2_idx
  ON Door (room2);

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
CREATE INDEX slot_room_idx
  ON Slot (room);

CREATE TABLE Interactive (
  id           SERIAL PRIMARY KEY,
  room         INT REFERENCES Room (id),
  name         VARCHAR(100),
  description  VARCHAR(500),
  isAccessible BOOLEAN,
  args         VARCHAR(100) --this is an argument string. Arguments has to be divided with ' '
);
CREATE INDEX interactive_room_idx
  ON Interactive (room);

CREATE TABLE Action (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(100),
  resultCode INT DEFAULT 0,
  resultMsg  VARCHAR(500)
);

CREATE TABLE Labyrinth (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) UNIQUE,
  description VARCHAR(500),
  startRoom   INT REFERENCES Room (id)
);
CREATE INDEX labyrinth_room_idx
  ON Labyrinth (startRoom);

CREATE TABLE InteractiveObjectNeed (
  id          SERIAL PRIMARY KEY,
  interactive INT REFERENCES Interactive (id),
  item        INT REFERENCES Item (id),
  UNIQUE (interactive, item)
);
CREATE INDEX interactiveObjectNeed_idx
  ON InteractiveObjectNeed (interactive, item);

CREATE TABLE ActionInteractiveLink (
  id          SERIAL PRIMARY KEY,
  action      INT REFERENCES Action (id),
  interactive INT REFERENCES Interactive (id),
  UNIQUE (action, interactive)
);
CREATE INDEX actionInteractiveLink_idx
  ON ActionInteractiveLink (action, interactive);

CREATE TABLE ActionInteractiveSwitch (
  id          SERIAL PRIMARY KEY,
  action      INT REFERENCES Action (id),
  interactive INT REFERENCES Interactive (id),
  newState    BOOLEAN,
  UNIQUE (action, interactive)
);
CREATE INDEX ActionInteractiveSwitch_idx
  ON ActionSlotSwitch (action, slot);

CREATE TABLE ActionDoorSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id),
  door     INT REFERENCES Door (id),
  newState BOOLEAN,
  UNIQUE (action, door)
);
CREATE INDEX actionDoorSwitch_idx
  ON ActionDoorSwitch (action, door);

CREATE TABLE ActionSlotSwitch (
  id       SERIAL PRIMARY KEY,
  action   INT REFERENCES Action (id),
  slot     INT REFERENCES Slot (id),
  newState BOOLEAN,
  UNIQUE (action, slot)
);
CREATE INDEX actionSlotSwitch_idx
  ON ActionSlotSwitch (action, slot);

CREATE TABLE SlotItemLink (
  id   SERIAL PRIMARY KEY,
  slot INT REFERENCES Slot (id),
  item INT REFERENCES Item (id),
  UNIQUE (slot, item)
);
CREATE INDEX SlotItemLink_idx
  ON SlotItemLink (slot, item);

CREATE TABLE LabyrinthRoomLink (
  id        SERIAL PRIMARY KEY,
  room      INT REFERENCES Room (id),
  labyrinth INT REFERENCES Labyrinth (id),
  UNIQUE (room, labyrinth)
);
CREATE INDEX labyrinthRoomLink_idx
  ON LabyrinthRoomLink (room, labyrinth);

CREATE TABLE LabyrinthActionLink (
  id        SERIAL PRIMARY KEY,
  action    INT REFERENCES Action (id),
  labyrinth INT REFERENCES Labyrinth (id),
  UNIQUE (action, labyrinth)
);
CREATE INDEX labyrinthActionLink_idx
  ON LabyrinthActionLink (action, labyrinth);
