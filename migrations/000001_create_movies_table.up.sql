CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    release_date date NOT NULL,
    runtime interval NOT NULL,
    genres text[] NOT NULL,     
    language text NOT NULL
);

CREATE TABLE IF NOT EXISTS cinemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    city text NOT NULL    
);


CREATE TABLE IF NOT EXISTS cinema_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    seats_count int NOT NULL,
    cinema_id UUID NOT NULL references cinemas(id)    
);

CREATE TABLE IF NOT EXISTS showtimes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL references movies(id),
    cinema_room_id UUID NOT NULL references cinema_rooms(id),
    show_start timestamptz NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS fare (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL references showtimes(id),
    fare_type_description text NOT NULL,
    price numeric NOT NULL 
);

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    surname text NOT NULL,
    phone text NOT NULL check(phone ~ '^\+[1-9]\d{6,14}$')
);


CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL references showtimes(id),
    customer_id UUID NOT NULL references customers(id),
    price numeric NOT NULL
);

CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL references showtimes(id),
    row_number int NOT NULL,
    seat_number int NOT NULL,
    seat_status text NOT NULL check(seat_status in ('available', 'booked')) DEFAULT 'available',
    fare_id UUID NOT NULL references fare(id),
    reservation_id UUID references reservations(id) 
);