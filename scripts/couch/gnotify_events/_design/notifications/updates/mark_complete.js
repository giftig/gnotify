function(doc, req) {
  var response = {
    headers: {
      'Content-Type': 'text/plain'
    },
    body: '',
    code: 200
  };

  if (!doc) {
    response.code = 404;
    return [null, response];
  }

  doc.complete = true;
  return [doc, response];
}
