'use strict';

const express = require('express');
const router = express.Router();

router.get('/', (req, res, next) => {
  res.render('pages/index', {title: 'Random Chat Demo'});
});

module.exports = router;
