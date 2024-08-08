const db = require('./config/db');

const initDb = async () => {
    const createTableQuery = `
    CREATE TABLE IF NOT EXISTS cmdb_assets (
      id INT AUTO_INCREMENT PRIMARY KEY,
      ip VARCHAR(15) NOT NULL,
      application_system VARCHAR(255),
      application_manager VARCHAR(255),
      overall_manager VARCHAR(255),
      is_virtual_machine BOOLEAN,
      resource_pool VARCHAR(255),
      data_center VARCHAR(255),
      rack_location VARCHAR(255),
      sn_number VARCHAR(255),
      out_of_band_ip VARCHAR(15),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )
  `;

    try {
        await db.query(createTableQuery);
        console.log('Database tables initialized successfully.');
    } catch (error) {
        console.error('Error initializing database tables:', error);
    }
};

module.exports = initDb;
