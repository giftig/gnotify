function(doc) {
  if (!doc.complete) {
    emit([doc.recipient, doc.time], null);
  }
}
