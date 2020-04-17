'use strict';

const express = require('express');
const router = express.Router();
const request = require('request');

router.get('/', (req, res, next) => {
  res.render('pages/index', {title: 'Random Chat Demo'});
});

router.post('/join', (req, res, next) => {
  let username = req.body.username;
  // request({
  //   url: 'http://localhost:9000/join',
  //   method: 'POST',
  //   headers: {
  //     'Accept': 'application/json',
  //     'Content-Type': 'application/json'
  //   },
  //   body: {username: username},
  //   json: true
  // }, (err, response, body) => {
  //   if(err || response.statusCode != '200'){
  //     console.log('error');
  //   }else{
  //       console.log(body)
  //   }
  // });

  var cookie = req.cookies.username;
  if (cookie === undefined) {
    // no: set a new cookie
    res.cookie('username', username, { maxAge: 900000, httpOnly: true });
  } else {
    // yes, cookie was already present 
    console.log('cookie exists', cookie);
  }
  res.redirect('/chat');
});

router.get('/chat', (req, res, next) => {
  var username = req.cookies.username;
  if (username === undefined) {
    return res.redirect('/');
  }
  res.render('pages/chat', {title: 'Random Chat Demo', username: username});
});

module.exports = router;
