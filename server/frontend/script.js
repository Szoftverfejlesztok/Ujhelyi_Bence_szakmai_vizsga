// init runs only once when the page loads
function initMain() {
    getDevices()
        .then(data => {
            if (data !== null) {
                const devices = JSON.parse(data);
                for (let i = 0; i < devices.length; i++) {
                    addDevice(i, devices[i].device);
                    if ( devices[i].state === true ) {
                        setSliderState(makeStringFancy(devices[i].device), devices[i].state)
                    }
                }
            }
        });
}

function initStatistics() {
    getDevicesUptime()
        .then(data => {
            if (data !== null) {
                const devices = JSON.parse(data);
                updateChart(devices);
            }
        });
}

// getDevices get the devices and their states from the database
function getDevices() {
    return fetch('/api/getDevices')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.text();
        })
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error(`Error making the GET request: ${error.message}`);
            return null;
        });
}

// makeStringFancy replace underscore with a space and make every first character uppercase
function makeStringFancy(str) {
    const words = str.split('_');
    const capitalizedWords = words.map(word => word.charAt(0).toUpperCase() + word.slice(1));
    return capitalizedWords.join(' ');
}

// addDevice add devices to the frontend, triggered only when the page loads
function addDevice(id, deviceName) {
    const cardDiv = document.createElement("div");
    cardDiv.className = "card";
    cardDiv.style.width = "auto";
    cardDiv.style.border = "0";

    const cardBody = document.createElement("div");
    cardBody.className = "card-body";
    cardBody.style.padding = "0 0 0 10px";

    const switchContainer = document.createElement("div");
    switchContainer.className = "switch-container";

    const switchLabel = document.createElement("label");
    switchLabel.className = "switch";
    switchLabel.innerHTML = "<input type=\"checkbox\"><span class=\"slider\"></span>";

    const deviceNameElement = document.createElement("h4");
    deviceNameElement.className = 'card-title';
    deviceNameElement.innerText = makeStringFancy(deviceName);

    switchContainer.appendChild(deviceNameElement);
    switchContainer.appendChild(switchLabel);

    cardBody.appendChild(switchContainer);
    cardDiv.appendChild(cardBody);

    document.getElementById("cardContainer").appendChild(cardDiv);

    const switchInput = cardDiv.querySelector('.switch input');
    switchInput.addEventListener('click', function() {
        sendNewState(deviceName, switchInput.checked);
    });
}

// setSliderState set the state of the given toggle switch based on the given value
function setSliderState(cardTitle, newState) {
    const cardTitles = document.querySelectorAll('.card-title');
    cardTitles.forEach(titleElement => {
        if (titleElement.textContent === cardTitle) {
            const slider = titleElement.nextElementSibling.querySelector('.switch input');
            if (slider) {
                slider.checked = newState;

                const changeEvent = new Event('change', { bubbles: true });
                slider.dispatchEvent(changeEvent);
            }
        }
    });
}

// sendNewState send the new state of a switch, triggered by a toggle switch
function sendNewState(name, state) {
    const data = {
        device: name,
        state: state
    };
    const jsonData = JSON.stringify(data);

    const requestOptions = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonData
    };

    fetch('/api/addRecord', requestOptions)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.json();
        })
        .then(data => {
            console.log(`Device ${makeStringFancy(data.device)} set to ${data.state}`);
        })
        .catch(error => {
            console.error('Error:', error);
        });

}

// getDevicesUptime get the devices and their uptime in order from the database
function getDevicesUptime() {
    return fetch('/api/getDevicesUptime')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.text();
        })
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error(`Error making the GET request: ${error.message}`);
            return null;
        });
}

// updateChart
function updateChart(devices) {
    let xValues = [];
    let yValues = [];

    for (let i = 0; i < devices.length; i++) {
        let uptime = precise(devices[i].uptime / 3600)
        let label = devices[i].device + " (" + uptime + "h)"

        xValues.push(makeStringFancy(label))
        yValues.push(uptime)
    }

    new Chart("uptimeChart", {
    type: "bar",
    data: {
        labels: xValues,
        datasets: [{
        backgroundColor: "blue",
        data: yValues
        }]
    },
    options: {
        legend: {display: false},
        title: {
        display: true,
        text: "Uptimes in the last 24h"
        }
    }
    });
}

function precise(x) {
    return x.toPrecision(2);
}
  