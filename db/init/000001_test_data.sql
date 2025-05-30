INSERT INTO public.users (id, prefered_user_name, given_name, family_name, email) VALUES
    ('74b19236-2243-4b2a-a71f-d24d5fa4c58a', 'owner', 'Colette', 'Tatou', 'owner@solei.com'),
    ('e2e7eb01-e704-4435-ab4e-93805b7593ed', 'customer', 'Anton', 'Ego', 'customer@solei.com');

INSERT INTO public.addresses (id, contry, city, street, phone, user_id) VALUES
    ('5a7acf66-7ee2-45e8-984d-64d5464cb263', 'France','Paris', 'Avenue Victor Hugo 11', '0031 979 887', '74b19236-2243-4b2a-a71f-d24d5fa4c58a'),
    ('30b1ac63-6a3f-4038-a010-5c71d2ed372d', 'France','Nice', 'Place Garibaldi 17', '0031 556 921', 'e2e7eb01-e704-4435-ab4e-93805b7593ed');

INSERT INTO public.categories (id, name) VALUES 
    ('6870aac5-38af-451d-846a-5572885aadd3', 'Starters'),
    ('147e8e67-613a-42be-a1a8-cce202e246f7', 'Entrees'),
    ('07509e26-8b5a-4e22-af32-1f85af73e8a0', 'Deserts');

INSERT INTO public.meals (id, name, description, cost, category_id) VALUES
    ('d98e2eba-e21b-4f6b-a599-60ab1c4c7237', 'Bruschetta', 'Toasted bread slices topped with fresh tomatoes, herbs, and sometimes cheese.', '5.10','6870aac5-38af-451d-846a-5572885aadd3'),
    ('b1c8f3d2-4e5a-4b0c-9f6d-2e7f8c1a2b3c', 'Caesar Salad', 'Crisp romaine lettuce with Caesar dressing, croutons, and Parmesan cheese.', '18.20','6870aac5-38af-451d-846a-5572885aadd3'),
    ('f3e2b1c4-5d6e-4f7a-b8c9-d0e1f2a3b4c5', 'Grilled Salmon', 'Salmon fillet grilled to perfection, served with seasonal vegetables.', '35.30','147e8e67-613a-42be-a1a8-cce202e246f7'),
    ('761c7614-54c1-47d6-afa1-7f678b633955', 'Chocolate Mousse', 'Rich and creamy chocolate mousse topped with whipped cream.', '12.40','07509e26-8b5a-4e22-af32-1f85af73e8a0');

INSERT INTO public.orders (id, price, status, user_id) VALUES
    ('6ec514e1-9111-4ac6-a327-f3421522f363', '5.10', 'DRAFT','e2e7eb01-e704-4435-ab4e-93805b7593ed');

INSERT INTO public.orderitems (id, amount, comment, meal_id, order_id, user_id) VALUES
    ('83eb039f-cc0a-4c31-9441-4f006cd5237e', '1', 'Less Mayo please!','d98e2eba-e21b-4f6b-a599-60ab1c4c7237', '6ec514e1-9111-4ac6-a327-f3421522f363','e2e7eb01-e704-4435-ab4e-93805b7593ed');
