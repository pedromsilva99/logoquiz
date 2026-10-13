-- +goose Up
INSERT INTO levels (number, title) VALUES (1, 'Marcas do dia a dia');

INSERT INTO logos (level_id, position, name, image_key)
SELECT l.id, v.position, v.name, v.image_key
FROM levels l,
     (VALUES
         (1,  'Continente',               'continente'),
         (2,  'Pingo Doce',               'pingo-doce'),
         (3,  'MEO',                      'meo'),
         (4,  'NOS',                      'nos'),
         (5,  'Caixa Geral de Depósitos', 'cgd'),
         (6,  'Millennium BCP',           'millennium-bcp'),
         (7,  'Galp',                     'galp'),
         (8,  'TAP Air Portugal',         'tap'),
         (9,  'CTT',                      'ctt'),
         (10, 'Super Bock',               'super-bock'),
         (11, 'Sagres',                   'sagres'),
         (12, 'Delta Cafés',              'delta')
     ) AS v (position, name, image_key)
WHERE l.number = 1;

INSERT INTO logo_aliases (logo_id, alias)
SELECT lg.id, v.alias
FROM logos lg
JOIN (VALUES
         ('cgd',            'CGD'),
         ('cgd',            'Caixa'),
         ('millennium-bcp', 'Millennium'),
         ('millennium-bcp', 'BCP'),
         ('galp',           'Galp Energia'),
         ('tap',            'TAP'),
         ('ctt',            'CTT Correios'),
         ('ctt',            'Correios de Portugal'),
         ('delta',          'Delta')
     ) AS v (image_key, alias) ON v.image_key = lg.image_key;

-- +goose Down
DELETE FROM levels WHERE number = 1;
