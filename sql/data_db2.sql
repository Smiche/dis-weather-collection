
INSERT INTO public.station_local(name, "number", organization, type, latitude, longitude, altitude, city)
VALUES ( 'central_station', 3000, 2, 'synoptical', 59, 18, 10, 'Stockholm');

INSERT INTO public.device_local(station, name)
VALUES (1, 'ATMOS 41');
