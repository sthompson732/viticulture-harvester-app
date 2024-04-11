-- Insert initial data into the vineyards table with bounding box coordinates
INSERT INTO vineyards (name, location, bbox) VALUES
('King Family Vineyards', 'Crozet, Virginia', 'POLYGON((-78.7014 38.0693, -78.7014 38.0650, -78.6971 38.0650, -78.6971 38.0693, -78.7014 38.0693))');

-- Insert initial data into the images table
-- Assuming vineyard_id for King Family Vineyards is 1
INSERT INTO images (vineyard_id, image_url, bbox, captured_at, description) VALUES
(1, 'https://storage.googleapis.com/{bucket_name}/king_family_vineyard1.jpg', 'POLYGON((-78.7014 38.0693, -78.7014 38.0680, -78.6990 38.0680, -78.6990 38.0693, -78.7014 38.0693))', '2023-10-01', 'Overview of King Family Vineyards'),
(1, 'https://storage.googleapis.com/{bucket_name}/king_family_vineyard2.jpg', 'POLYGON((-78.7000 38.0670, -78.7000 38.0655, -78.6975 38.0655, -78.6975 38.0670, -78.7000 38.0670))', '2023-10-02', 'Detailed shot of Merlot section');

-- Insert initial data into the satellite_imagery table with GCS URLs
INSERT INTO satellite_imagery (vineyard_id, image_url, bbox, captured_at, resolution) VALUES
(1, 'https://storage.googleapis.com/{bucket_name}/sat_king_family1.jpg', 'POLYGON((-78.7014 38.0693, -78.7014 38.0650, -78.6971 38.0650, -78.6971 38.0693, -78.7014 38.0693))', '2023-09-15', 0.1),
(1, 'https://storage.googleapis.com/{bucket_name}/sat_king_family2.jpg', 'POLYGON((-78.7005 38.0688, -78.7005 38.0665, -78.6980 38.0665, -78.6980 38.0688, -78.7005 38.0688))', '2023-09-16', 0.1);

-- Insert initial data into the soil_data table with realistic data points
INSERT INTO soil_data (vineyard_id, data, location, sampled_at) VALUES
(1, '{"ph": 6.7, "nutrients": {"nitrogen": 45, "phosphorus": 15, "potassium": 20}, "organic_matter": 5.6}', 'POINT(-78.6995 38.0675)', '2024-03-20');