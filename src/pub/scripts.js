function sendData(pcId) {
    fetch(`http://turret.okke.dev/select/${pcId}`, {
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
        .catch(_ => {
            return
        });
}
