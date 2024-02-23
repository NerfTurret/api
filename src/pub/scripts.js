function sendData(pcId) {
    fetch(`http://127.0.0.1:3000/select/${pcId}`, {
        method: 'GET',
      })
        .then(response => {
          if (!response.ok) {
              console.error(response)
            throw new Error('Network response was not ok');
          }
          return response.json();
        })
        .then(data => {
          console.log(data.json());
        })
        .catch(error => {
          console.error('There was a problem with the fetch operation:', error.json());
        });
}
