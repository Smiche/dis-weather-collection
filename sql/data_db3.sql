
INSERT INTO public.station_local(name, "number", organization, type, latitude, longitude, altitude, city)
VALUES ( 'west_station', 4000, 3, 'aeronautical', 59, 10, 7, 'Oslo');

INSERT INTO public.device_local(station, name)
VALUES (1, 'Vaisala WXT534');
