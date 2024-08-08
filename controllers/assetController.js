const assetModel = require('../models/assetModel');
const xlsx = require('xlsx'); // 确保导入 xlsx 模块

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
        console.log('Request received:', req); // 检查请求内容
        if (!req.file) {
            console.log('No file found in request'); // 文件不存在时的日志
            throw new Error('No file uploaded');
        }

        console.log('File received:', req.file); // 检查文件内容
        const workbook = xlsx.read(req.file.buffer, { type: 'buffer' });
        console.log('Workbook parsed:', workbook); // 检查工作簿内容

        const sheetName = workbook.SheetNames[0];
        console.log('Sheet name:', sheetName); // 检查工作表名称

        const sheet = workbook.Sheets[sheetName];
        const data = xlsx.utils.sheet_to_json(sheet);
        console.log('Data parsed from sheet:', data); // 检查解析的数据

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
