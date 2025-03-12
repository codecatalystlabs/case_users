
const ENDPOINTS = {'create-allergy': '/clients/view/allergy/create',
                   'create-illness': '/clients/view/chronic/create',
                   'create-contact': '/clients/view/contact/create',
                   'create-identifier': '/clients/view/identifier/create',

                   'delete-allergy': '/clients/view/allergy/delete',
                   'delete-illness': '/clients/view/chronic/delete',
                   'delete-contact': '/clients/view/contact/delete',
                   'delete-identifier': '/clients/view/identifier/delete',
                
                   'queue-patient': '/clients/view/add-to-queue',
                   'unqueue-patient': '/clients/view/remove-from-queue',
                
                   'create-complaint':'/clients/view/complaint/create',
                   'update-complaint':'/clients/view/complaint/update',
                   'delete-complaint':'/clients/view/complaint/delete',

                   'create-measure':'/clients/view/measure/create',
                   'update-measure':'/clients/view/measure/update',
                   'delete-measure':'/clients/view/measure/delete',
                }


const REDIRECT_URLS = {'patient-view': '/clients/view/'}

const CLINIC_ENDPOINTS = {
    'fetch-clinic-charges': '/facility/clinic/charge',
    'fetch-physicians-by-facility':'/facility/physicians'
}

let selected_services = [];


function showLoading(){
    const modal = document.getElementById('loadingBackdrop')
    modal.style.display = 'block';
    //$("#loadingBackdrop").modal('show');
}

function hideLoading(){
    //console.log("hiding loading")
    const modal = document.getElementById('loadingBackdrop')
    modal.style.display = 'none';
    //$("#loadingBackdrop").modal('hide');
}

function getCsrfToken(){
    return document.getElementById("csrf_token").value
}

function getPatientID(){
    const patient_id = document.getElementById("id").value
    if(patient_id !== null){
        return patient_id
    }
    return ""
}


function displayMessage(feedback, error){
    /* error value is 1 or 0 */
    let bgColor;
    let fontColor;
    let borderColor;
    if(error){
        //bg_gradient = "linear-gradient(to right, #e04e19, #e04e19)"
        borderColor = "1px solid #E32230";
        bgColor = "#e6e7ee";
        fontColor = "#E32230";
    }
    else{
        borderColor = "1px solid #99BBE8";
        bgColor = "#e6e7ee";
        fontColor = "#014f86";
    }
    Toastify({
        text: feedback,
        duration: 5000,
        newWindow: false,
        close: true,
        gravity: "top", // `top` or `bottom`
        position: "center", // `left`, `center` or `right`
        stopOnFocus: true, // Prevents dismissing of toast on hover
        style: {
            border: borderColor,
            background: bgColor,
            color: fontColor
        },
    }).showToast();
}


function showFeedback(message, error_status){
    // error_status is 1 or 0
    displayMessage(message, error_status)
}


function updateSelectedAllergy(value){
    const selected_allergy_input = document.getElementById("selected_allergy_id")
    selected_allergy_input.value = value
}

function updateSelectedIllness(value){
    const selected_illness_input = document.getElementById("selected_illness_id")
    selected_illness_input.value = value
}

function updateSelectedContact(value){
    const selected_contact_input = document.getElementById("selected_contact_id")
    selected_contact_input.value = value
}

function updateSelectedIdentifier(value){
    const selected_identifier_input = document.getElementById("selected_identifier_id")
    selected_identifier_input.value = value
}

function highlightTableRowAllergy(row) {
    // Remove highlight from previously selected row, if any
    const previouslySelectedRow = document.querySelector('tr.selected');
    if (previouslySelectedRow) {
    previouslySelectedRow.classList.remove('selected');
    }
   // Highlight the clicked row
    row.classList.add('selected');
    const dataValue = row.getAttribute('data-value');
    updateSelectedAllergy(dataValue)
   }


function highlightTableRowIllness(row) {
    const previouslySelectedRow = document.querySelector('tr.selected');
    if (previouslySelectedRow) {
    previouslySelectedRow.classList.remove('selected');
    }
    row.classList.add('selected');
    const dataValue = row.getAttribute('data-value');
    updateSelectedIllness(dataValue)
   }


function highlightTableRowContact(row) {
    const previouslySelectedRow = document.querySelector('tr.selected');
    if (previouslySelectedRow) {
    previouslySelectedRow.classList.remove('selected');
    }
    row.classList.add('selected');
    const dataValue = row.getAttribute('data-value');
    updateSelectedContact(dataValue)
   }


function highlightTableRowIdentifier(row) {
    const previouslySelectedRow = document.querySelector('tr.selected');
    if (previouslySelectedRow) {
    previouslySelectedRow.classList.remove('selected');
    }
    row.classList.add('selected');
    const dataValue = row.getAttribute('data-value');
    updateSelectedIdentifier(dataValue)
   }


function getClinicIDs(){
    const patient_id = document.getElementById("patient_id").value
    if(patient_id !== null){
        return patient_id
    }
    return ""
}


function serviceExists(array, service_id) {
    return array.includes(service_id);
  }

function addToServiceList(service_id) {
    if (!serviceExists(selected_services, service_id)){
        selected_services.push(service_id);
    }
}

function removeFromServiceList(service_id) {
    if (serviceExists(selected_services, service_id)){
        selected_services = selected_services.filter(item => item !== service_id)
    }
}


function handleSelection(checkbox_id) {
    var checkBox = document.getElementById(checkbox_id);
    const clinic_id = checkBox.getAttribute('data-clinic');
    const service_id = checkBox.getAttribute('data-service');
    if (checkBox.checked == true){
        fetch_charges(clinic_id, service_id)
        addToServiceList(service_id)
    } else {
        removeCharges(clinic_id, service_id)
        removeFromServiceList(service_id)
    }
}


function addCharges(data, clinic_id, service_id){
    const row_div_id = `${clinic_id}_${service_id}` //make unique id for row based on clinic and service identifiers
    const table_row = `
        <tr id="${row_div_id}">
            <td class="td-basic"> ${data.service}</td>
            <td class="td-basic"> ${data.cost}</td>
            <td class="td-basic"> ${data.optional}</td>
        </tr>`
    $('#charges_table_body').append(table_row);
}


function removeCharges(clinic_id, service_id){
    table_div_id = `${clinic_id}_${service_id}`;
    const element = document.getElementById(table_div_id)
    element.remove();
}


/* function expects an array of physicians */
function updatePhysicianDropdown(physicians){
    var data = ``
    for (var index in physicians){
        data += `<option value=${physicians[index]}> ${physicians[index]}</option>`
    }
    $("#select_physician").empty().html(data);
}

/* perform ajax call */

function submit_ajax(section){

    var token = getCsrfToken()
    var end_point = "/clients/subs/create"
    var patient_id = getPatientID()

    var filter_s = '[id^="' + section + '-"]'
    var elements = document.querySelectorAll( filter_s );

    var post_data = {
        'client':               patient_id,
        'whatisthis':           section,
    };

    elements.forEach(function(element) {
        var dashPosition = element.id.indexOf('-');
        var propertyName = element.id.slice(dashPosition); 
        post_data[propertyName] = element.value;
    });

    $.ajax({
        url: end_point,
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000, //10 seconds
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
}

function delete_ajax(section){

    var token = getCsrfToken()
    var end_point = "/clients/subs/delete"
    var patient_id = getPatientID()

    var filter_s = '[id^="' + section + '-id"]'
    var elements = document.querySelectorAll( filter_s );

    var post_data = {
        'client':               patient_id,
        'whatisthis':           section,
    };

    elements.forEach(function(element) {
        var dashPosition = element.id.indexOf('-');
        var propertyName = element.id.slice(dashPosition); 
        post_data[propertyName] = element.value;
    });

    $.ajax({
        url: end_point,
        type: "DELETE",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000, //10 seconds
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
}

/** Queue to Clinic **/
function post_queue(){
    var token = getCsrfToken()
    var patient_id = getPatientID()
    post_data = {
        'patient_id': patient_id,
        'services_list': selected_services,
        'medical_scheme': document.getElementById('select_scheme').value,
        'insurance_code': document.getElementById('insurance_code').value,
        'waiver': document.getElementById('select_waiver').value,
        'pricing_package': document.getElementById('select_package').value,
        'facility': document.getElementById('select_facility').value,
        'physician': document.getElementById('select_physician').value,
    }

    $.ajax({
        url: ENDPOINTS['queue-patient'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {'X-CSRFToken': token, 'csrftoken':token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
                //window.location.replace(REDIRECT_URLS['patient-view']);
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
};

// Remove from Queue
function delete_from_queue(){
    var token = getCsrfToken()
    var patient_id = getPatientID()
    post_data = {
        'patient_id': patient_id,
        'clinic_id': document.getElementById('selected_clinic_id').value,
    }
    
    $.ajax({
        url: ENDPOINTS['unqueue-patient'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {'X-CSRFToken': token, 'csrftoken':token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
                //window.location.replace(REDIRECT_URLS['patient-view']);
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
};



/** Fetch Clinic Charge **/
function fetch_charges(clinic_id, service_id){
    var token = getCsrfToken()
    var post_data = {
        'service_id': service_id,
        'clinic_id': clinic_id
    }
    $.ajax({
        url: CLINIC_ENDPOINTS['fetch-clinic-charges'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                //console.log("RESP ", response_data)
                addCharges(response_data, clinic_id, service_id)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            console.log(response_data)
            showFeedback(response_data.feedback, 1)
        }
    });
};


/** Get Physicians by Facility **/
function get_physicians(){
    var token = getCsrfToken()
    var facility_id = document.getElementById("select_facility").value
    var post_data = {
        'facility_id': facility_id
    }
    $.ajax({
        url: CLINIC_ENDPOINTS['fetch-physicians-by-facility'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                updatePhysicianDropdown(response_data.physicians_list)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
};



function updateComplaint(complaint_id, complaint_field, new_value){
    // complaint field can be the name or its description
    var token = getCsrfToken()
    var post_data = {
        'complaint_id': complaint_id,
        'complaint_field': complaint_field,
        'edited_complaint_value': new_value
    }
    $.ajax({
        url: ENDPOINTS['update-complaint'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
}


document.addEventListener('DOMContentLoaded', function() {
    function makeComplaintEditable(element) {
        element.addEventListener('click', function() {
            const currentText = this.textContent;
            const element_id = this.id
            const {_, complaint_field, complaint_id} = element_id.split("_")
            const input = document.createElement('input');
            input.type = 'text';
            input.value = currentText;
            input.className = 'editable-input';

            this.parentNode.replaceChild(input, this);
            input.focus();

            //Add onChange listener 
            input.addEventListener('change', function() {
                //console.log(complaint_id, complaint_field, this.value);
                updateComplaint(complaint_id, complaint_field, this.value)
            });

            input.addEventListener('blur', function() {
                const newText = this.value;
                const span = document.createElement('span');
                span.className = 'editable editable-complaint';
                span.id = element_id;
                span.textContent = newText;

                this.parentNode.replaceChild(span, this);

                // Reattach the click event
                makeComplaintEditable(span);
            });

            input.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    this.blur();
                }
            });
        });
    }
    const editableElements = document.querySelectorAll('.editable-complaint');
    editableElements.forEach(makeComplaintEditable);
});



function updateMeasure(measure_id, measure_field, new_value){
    // measure field can be the attribute(height, weight), its unit (metres, kg) or the value(200cm, 30kg)
    var token = getCsrfToken()
    var post_data = {
        'measure_id': measure_id,
        'measure_field': measure_field,
        'edited_measure_value': new_value
    }
    $.ajax({
        url: ENDPOINTS['update-measure'],
        type: "POST",
        data: JSON.stringify(post_data),
        headers: {
            "X-CSRFToken":token, 
            "csrftoken":token},
        contentType: "application/json",
        processData: false,
        cache: false,
        timeout: 10000,
        beforeSend: showLoading(),
        success: function(response_data){
            hideLoading()
            if(response_data.message == 'OK'){
                showFeedback(response_data.feedback, 0)
            }
            else{
                showFeedback(response_data.feedback, 1)
            }
        },
        error: function(response_data){
            hideLoading()
            showFeedback(response_data.feedback, 1)
        }
    });
}


document.addEventListener('DOMContentLoaded', function() {
    function makeMeasureEditable(element) {
        element.addEventListener('click', function() {
            const currentText = this.textContent;
            const element_id = this.id
            const {_, measure_field, measure_id} = element_id.split("_")
            const input = document.createElement('input');
            input.type = 'text';
            input.value = currentText;
            input.className = 'editable-input';

            this.parentNode.replaceChild(input, this);
            input.focus();

            //Add onChange listener 
            input.addEventListener('change', function() {
                //console.log(measure_id, measure_field, this.value);
                updateMeasure(measure_id, measure_field, this.value)
            });

            input.addEventListener('blur', function() {
                const newText = this.value;
                const span = document.createElement('span');
                span.className = 'editable editable-measure';
                span.id = element_id;
                span.textContent = newText;

                this.parentNode.replaceChild(span, this);

                // Reattach the click event
                makeMeasureEditable(span);
            });

            input.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    this.blur();
                }
            });
        });
    }
    const editableElements = document.querySelectorAll('.editable-measure');
    editableElements.forEach(makeMeasureEditable);
});


// Function to update the selected value of the dropdown
function updateDropdownValue(newValue) {
    var dropdown = document.getElementById('myDropdown');

    if (!dropdown) {
        console.log('Dropdown element not found');
        return;
    }

    var options = Array.from(dropdown.options);
    var validValue = options.some(option => option.value === newValue);

    if (!validValue) {
        console.log('Invalid value provided');
        return;
    }

    dropdown.value = newValue;
}

// Example: Update the dropdown to select "option2"
updateDropdownValue('option2');


function readDataValues(element){
    const dataValues = element.getAttribute('data-values')
    var result = dataValues.split("|") // gives us an array like ['the_attribute', 'the_unit', 'the_value] i.e ['height', 'feet', '21']
    return result
}

function updateMeasurementDropdownValues(attribute_data, unit_data, value_data){
    var attribute_dropdown = document.getElementById('select_attribute');
    var unit_dropdown = document.getElementById('select_unit');
    var value_input = document.getElementById('measure_value');

    attribute_dropdown.value = attribute_data
    unit_dropdown.value = unit_data
    value_input.value = value_data
}

function updateComplaintDropdownValues(complaint_name, description){
    var complaint_dropdown = document.getElementById('select_complaint');
    var value_input = document.getElementById('complaint_description');

    complaint_dropdown.value = complaint_name
    value_input.value = description
}

function updateStabilizationDropdownValues(stabilization_name, description){
    var stabilization_dropdown = document.getElementById('select_stabilization');
    var value_input = document.getElementById('stabilization_description');

    stabilization_dropdown.value = stabilization_name
    value_input.value = description
}

function showEditingTableRow(input_trow_id){
    var input_row = document.getElementById(input_trow_id);
    input_row.style.display = 'table-row';
}

function addMeasure(){
    showEditingTableRow('measure_row_input')
}

function editMeasure(element){
    const dataValues = element.getAttribute('data-values')
    const data = readDataValues(element);
    updateMeasurementDropdownValues(data[0], data[1], data[2])
    showEditingTableRow('measure_row_input')
}

function addComplaint(){
    showEditingTableRow('complaint_row_input')
}

function editComplaint(element){
    const data = readDataValues(element);
    updateComplaintDropdownValues(data[0], data[1])
    showEditingTableRow('complaint_row_input')
}


function addStabilization(){
    showEditingTableRow('stabilization_row_input')
}


function editStabilization(element){
    const data = readDataValues(element);
    updateStabilizationDropdownValues(data[0], data[1])
    showEditingTableRow('stabilization_row_input')
}

function deleteRow(row_id){
    var trow = document.getElementById(row_id)
    var table_body = trow.parentNode;
    table_body.removeChild(trow)
}
