-- Insert initial data into the vineyards table with bounding box coordinates
INSERT INTO vineyards (name, location, bbox) VALUES
('King Family Vineyards', 'Crozet, Virginia', 'POLYGON((-78.7014 38.0693, -78.7014 38.0650, -78.6971 38.0650, -78.6971 38.0693, -78.7014 38.0693))');

-- Insert initial data into the images table
-- Assuming vineyard_id for King Family Vineyards is 1
INSERT INTO images (vineyard_id, image_url, bbox, captured_at, description) VALUES
(1, 'https://storage.googleapis.com/{bucket_name}/king_family_vineyard1.jpg', 'POLYGON((-78.7014 38.0693, -78.7014 38.0680, -78.6990 38.0680, -78.6990 38.0693, -78.7014 38.0693))', '2024-01-01', 'Overview of King Family Vineyards'),
(1, 'https://storage.googleapis.com/{bucket_name}/king_family_vineyard2.jpg', 'POLYGON((-78.7000 38.0670, -78.7000 38.0655, -78.6975 38.0655, -78.6975 38.0670, -78.7000 38.0670))', '2024-01-02', 'Detailed shot of Merlot section');

-- Insert initial data into the satellite_imagery table with GCS URLs
INSERT INTO satellite_imagery (vineyard_id, image_url, bbox, captured_at, resolution) VALUES
(1, 'https://storage.googleapis.com/{bucket_name}/sat_king_family1.jpg', 'POLYGON((-78.7014 38.0693, -78.7014 38.0650, -78.6971 38.0650, -78.6971 38.0693, -78.7014 38.0693))', '2024-01-15', 0.1),
(1, 'https://storage.googleapis.com/{bucket_name}/sat_king_family2.jpg', 'POLYGON((-78.7005 38.0688, -78.7005 38.0665, -78.6980 38.0665, -78.6980 38.0688, -78.7005 38.0688))', '2024-01-16', 0.1);

-- Insert initial data into the soil_data table with multiple entries over different dates
INSERT INTO soil_data (vineyard_id, data, location, sampled_at) VALUES
(1, '{"ph": 6.7, "nutrients": {"nitrogen": 45, "phosphorus": 15, "potassium": 20}, "organic_matter": 5.6}', ST_SetSRID(ST_MakePoint(-78.6995, 38.0675), 4326), '2024-01-20'),
(1, '{"ph": 6.5, "nutrients": {"nitrogen": 50, "phosphorus": 18, "potassium": 25}, "organic_matter": 3.2}', ST_SetSRID(ST_MakePoint(-78.6995, 38.0675), 4326), '2024-02-20'),
(1, '{"ph": 6.8, "nutrients": {"nitrogen": 40, "phosphorus": 20, "potassium": 22}, "organic_matter": 4.0}', ST_SetSRID(ST_MakePoint(-78.6995, 38.0675), 4326), '2024-03-20');

-- Insert Pest Data with multiple observations
INSERT INTO pest_data (vineyard_id, description, observation_date, location, pest_type, severity) VALUES
(1, 'Detected mild aphid infestation on younger vines.', '2024-03-15 10:00:00+00', ST_SetSRID(ST_MakePoint(-78.7006, 38.0685), 4326), 'Aphids', 'Mild'),
(1, 'Observed increased ladybug activity, beneficial for aphid control.', '2024-03-20 09:30:00+00', ST_SetSRID(ST_MakePoint(-78.7008, 38.0687), 4326), 'Ladybugs', 'Beneficial'),
(1, 'Spotted severe grapevine moth larvae in the northeast section.', '2024-03-22 12:45:00+00', ST_SetSRID(ST_MakePoint(-78.6998, 38.0689), 4326), 'Grapevine Moth Larvae', 'Severe'),
(1, 'Mildew detected on several older vines, potentially harmful.', '2024-03-25 15:30:00+00', ST_SetSRID(ST_MakePoint(-78.7000, 38.0691), 4326), 'Mildew', 'Moderate'),
(1, 'Early signs of botrytis fungus observed on some vines.', '2024-03-28 16:00:00+00', ST_SetSRID(ST_MakePoint(-78.7002, 38.0683), 4326), 'Botrytis Fungus', 'Early stage');

-- Insert Weather Data with multiple entries for better pattern analysis
INSERT INTO weather_data (vineyard_id, temperature, humidity, observation_time, location) VALUES
(1, 22.5, 78.9, '2024-01-15 08:00:00+00', ST_SetSRID(ST_MakePoint(-78.7006, 38.0685), 4326)),
(1, 18.0, 82.0, '2024-02-15 08:00:00+00', ST_SetSRID(ST_MakePoint(-78.7006, 38.0685), 4326)),
(1, 19.5, 80.5, '2024-03-15 08:00:00+00', ST_SetSRID(ST_MakePoint(-78.7006, 38.0685), 4326));