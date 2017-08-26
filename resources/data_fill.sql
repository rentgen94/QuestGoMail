
GRANT ALL ON users TO go_user;
GRANT ALL ON users_id_seq TO go_user;
GRANT ALL ON Room TO go_user;
GRANT ALL ON Door TO go_user;
GRANT ALL ON Item TO go_user;
GRANT ALL ON Slot TO go_user;
GRANT ALL ON Interactive TO go_user;
GRANT ALL ON Action TO go_user;
GRANT ALL ON Labyrinth TO go_user;
GRANT ALL ON InteractiveObjectNeed TO go_user;
GRANT ALL ON ActionInteractiveLink TO go_user;
GRANT ALL ON ActionInteractiveSwitch TO go_user;
GRANT ALL ON ActionDoorSwitch TO go_user;
GRANT ALL ON ActionSlotSwitch TO go_user;
GRANT ALL ON SlotItemLink TO go_user;
GRANT ALL ON LabyrinthRoomLink TO go_user;
GRANT ALL ON LabyrinthActionLink TO go_user;

INSERT INTO Item (name, description, size) VALUES ('Axe', 'big', 100);
INSERT INTO Item (name, description, size) VALUES ('Sword', 'small', 200);
INSERT INTO Room (name, description) VALUES ('Black room', 'You are alone in the dark. Trying to touch everything around');
INSERT INTO Room (name, description) VALUES ('White room', 'Room is empty, but all bulbs are shining bright');
INSERT INTO Slot (room, name, capacity, isAccessible) VALUES (1, 'box', 1000, TRUE);
INSERT INTO Interactive (room, name, description, isAccessible, args) VALUES (1, 'button', 'descr', TRUE , 'alpha beta');
INSERT INTO Interactive (room, name, description, isAccessible, args) VALUES (1, 'wheel', 'strange wheel', FALSE, '');
INSERT INTO SlotItemLink (slot, item) VALUES (1, 1);
INSERT INTO SlotItemLink (slot, item) VALUES (1, 2);
INSERT INTO Labyrinth (name, description, startroom) VALUES ('lab', 'lab', 1);
INSERT INTO LabyrinthRoomLink (room, labyrinth) VALUES (1, 1);
INSERT INTO LabyrinthRoomLink (room, labyrinth) VALUES (2, 1);
INSERT INTO Action (name, resultCode, resultMsg) VALUES ('activate_wheel', 0, 'strange wheel activated');
INSERT INTO ActionInteractiveSwitch (action, interactive, newstate) VALUES (1, 2, TRUE);
INSERT INTO ActionInteractiveLink (action, interactive) VALUES (1, 1);
INSERT INTO LabyrinthActionLink (action, labyrinth) VALUES (1, 1);
INSERT INTO Door (room1, room2, name, isAccessible) VALUES (1, 2, 'Grey door', TRUE );
INSERT INTO InteractiveObjectNeed (interactive, item) VALUES (2, 1);