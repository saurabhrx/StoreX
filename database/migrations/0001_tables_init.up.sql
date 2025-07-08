CREATE TYPE role_type AS ENUM(
    'admin',
    'employee',
    'asset_manager',
    'employee_manager'
    );
CREATE TYPE empl_type AS ENUM(
    'full_time',
    'intern',
    'freelancer'
    );

CREATE TYPE owned_type AS ENUM(
    'remote_state',
    'client'
    );

CREATE TYPE status_type AS ENUM(
    'available',
    'assigned',
    'waiting_for_repair',
    'service',
    'damaged',
    'deleted'
    );

CREATE TYPE assets_type AS ENUM(
    'laptop',
    'mobile',
    'mouse',
    'monitor',
    'hard_disk',
    'pen_drive',
    'sim',
    'accessories'
    );

CREATE TABLE IF NOT EXISTS employees(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL ,
    last_name TEXT,
    email TEXT NOT NULL ,
    phone_no TEXT,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS active_user_idx ON employees(TRIM(LOWER(email))) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS employee_type(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID REFERENCES employees(id) NOT NULL ,
    type empl_type NOT NULL DEFAULT 'full_time',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS employee_role(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID REFERENCES employees(id) NOT NULL ,
    role role_type NOT NULL DEFAULT 'employee',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS assets(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand TEXT NOT NULL ,
    model TEXT NOT NULL ,
    serial_no TEXT NOT NULL ,
    asset_type assets_type NOT NULL ,
    owned_by owned_type NOT NULL DEFAULT 'remote_state',
    purchased_at DATE NOT NULL ,
    price NUMERIC(10,2) NOT NULL ,
    status status_type NOT NULL DEFAULT 'available',
    created_by UUID REFERENCES employees(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS active_assets_idx ON assets(serial_no) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS assigned_asset(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    employee_id UUID REFERENCES employees(id) NOT NULL ,
    start_date DATE NOT NULL ,
    end_date DATE NULL,
    reason_to_return TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS assigned_idx ON assigned_asset(asset_id,employee_id,start_date) ;

CREATE TABLE IF NOT EXISTS vendors(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL ,
    phone_no TEXT NOT NULL ,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS active_vendor_idx ON vendors(name,phone_no) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS warranty(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    start_date DATE,
    end_date DATE
);

CREATE TABLE IF NOT EXISTS services(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    vendor_id UUID REFERENCES vendors(id) NOT NULL,
    start_date DATE NOT NULL ,
    end_date DATE NOT NULL ,
    cost NUMERIC(10,2) NOT NULL ,
    remark TEXT
);

-- different devices--

CREATE TABLE IF NOT EXISTS laptop_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    ram INTEGER NOT NULL ,
    storage_capacity INTEGER NOT NULL ,
    processor TEXT NOT NULL ,
    os TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS laptop_idx ON laptop_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS mobile_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    ram INTEGER NOT NULL,
    storage_capacity INTEGER NOT NULL ,
    os TEXT NOT NULL ,
    imei_1 TEXT NOT NULL ,
    imei_2 TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS mobile_idx ON mobile_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS monitor_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL ,
    screen_size DECIMAL(4,1) NOT NULL ,
    resolution TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS monitor_idx ON monitor_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS mouse_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    connection_type TEXT NOT NULL ,
    dpi INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS mouse_idx ON mouse_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS hard_disk_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    type TEXT NOT NULL ,
    capacity INTEGER NOT NULL,
    interface TEXT NOT NULL ,
    rpm INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS hard_drive_idx ON hard_disk_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS pen_drive_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    capacity INTEGER NOT NULL,
    interface TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS pen_drive_idx ON pen_drive_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS sim_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    sim_number TEXT NOT NULL ,
    career TEXT NOT NULL ,
    plan_type TEXT NOT NULL ,
    activation_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS sim_idx ON sim_specs(asset_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS accessories_specs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS accessories_idx ON accessories_specs(asset_id) WHERE archived_at IS NULL;