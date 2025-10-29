-- Seed para tests e2e

-- Insertar cliente de prueba
INSERT INTO client (id, name, email) VALUES
  ('11111111-1111-1111-1111-111111111111', 'Test Client', 'test@example.com');

-- Insertar short_link de prueba
INSERT INTO short_link (id, original_route, client_id, code) VALUES
  ('22222222-2222-2222-2222-222222222222', 'https://example.com', '11111111-1111-1111-1111-111111111111', 'testcode');

-- Insertar analíticas de prueba
INSERT INTO link_analytics (id, link_id, total_views, unique_visitors) VALUES
  ('33333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222', 10, 5);

-- Insertar contador de vistas por país
INSERT INTO link_country_view_counter (id, link_id, country_code, total_views, unique_visitors) VALUES
  ('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 'ES', 7, 3);

-- Insertar visita de prueba
INSERT INTO link_visit (id, link_id, ip, user_agent) VALUES
  ('55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', '127.0.0.1', 'test-agent');
