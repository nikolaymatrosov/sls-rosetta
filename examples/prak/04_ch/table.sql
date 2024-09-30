CREATE TABLE weather
(  LocalDateTime DateTime,
   LocalDate Date,
   Month Int8,
   Day Int8,
   TempC Float32,
   Pressure Float32,
   RelHumidity Int32,
   WindSpeed10MinAvg Int32,
   VisibilityKm Float32,
   City String
) ENGINE=MergeTree
    ORDER BY LocalDateTime;

SELECT * FROM weather;