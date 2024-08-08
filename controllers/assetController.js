const assetModel = require('../models/assetModel');

const getAssets = async (req, res) => {
    const { ip, sn } = req.query;
    const field = ip ? 'ip' : 'sn_number';
    const value = ip || sn;
    try {
        const assets = await assetModel.getAssetsByField(field, value);
        res.json(assets);
    } catch (error) {
        res.status(500).json({ error: 'Database query error' });
    }
};

const createAsset = async (req, res) => {
    try {
        const assetData = req.body;
        const insertId = await assetModel.addAsset(assetData);
        res.status(201).json({ id: insertId });
    } catch (error) {
        console.error('Database insert error:', error);
        res.status(500).json({ error: 'Database insert error', details: error.message });
    }
};

const updateAsset = async (req, res) => {
    const { id } = req.params;
    try {
        await assetModel.updateAsset(id, req.body);
        res.json({ message: 'Asset updated successfully' });
    } catch (error) {
        res.status(500).json({ error: 'Database update error' });
    }
};

const deleteAsset = async (req, res) => {
    const { id } = req.params;
    try {
        await assetModel.deleteAsset(id);
        res.json({ message: 'Asset deleted successfully' });
    } catch (error) {
        res.status(500).json({ error: 'Database delete error' });
    }
};

const importAssets = async (req, res) => {
    try {
        const workbook = xlsx.read(req.file.buffer, { type: 'buffer' });
        const sheetName = workbook.SheetNames[0];
        const sheet = workbook.Sheets[sheetName];
        const data = xlsx.utils.sheet_to_json(sheet);

        for (const asset of data) {
            await assetModel.addAsset(asset);
        }

        res.json({ message: 'Data imported successfully' });
    } catch (error) {
        console.error('Error importing data:', error);
        res.status(500).json({ error: 'Failed to import data' });
    }
};

module.exports = {
    getAssets,
    createAsset,
    updateAsset,
    importAssets,
    deleteAsset
};
