INSERT INTO public.organization(
 name, country)
VALUES ( 'NOWeather', 'Norway');

INSERT INTO public.organization(
 name, country)
VALUES ( 'ENVIry', 'Finland');

INSERT INTO public.organization(
 name, country)
VALUES ( 'MeteoNat', 'Sweden');

INSERT INTO public.station_local(name, "number", organization, type, latitude, longitude, altitude, city)
VALUES ( 'south_station', 2000, 1, 'aeronautical', 22, 66, 10, 'Lappeenranta');

INSERT INTO public.device_local(station, name)
VALUES (1, 'Vaisala WXT536');