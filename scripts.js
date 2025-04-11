const API_URL = "http://localhost:8080/incidentes";

// Listar incidentes
async function fetchIncidents() {
    const response = await fetch(API_URL);
    const incidents = await response.json();
    const incidentList = document.getElementById("incident-list");
    incidentList.innerHTML = "";
    incidents.forEach(incident => {
        const div = document.createElement("div");
        div.className = "incident";
        div.innerHTML = `
            <h3>ID: ${incident.id}</h3>
            <p><strong>Empleado:</strong> ${incident.empleado}</p>
            <p><strong>Tipo de Equipo:</strong> ${incident.tipo_equipo}</p>
            <p><strong>Detalle del Problema:</strong> ${incident.detalle_problema}</p>
            <p><strong>Día del Problema:</strong> ${incident.dia_problema}</p>
            <p><strong>Estado:</strong> ${incident.estado}</p>
        `;
        incidentList.appendChild(div);
    });
}

// Crear incidente
async function createIncident() {
    const empleado = document.getElementById("empleado").value;
    const tipoEquipo = document.getElementById("tipo_equipo").value;
    const detalleProblema = document.getElementById("detalle_problema").value;
    const diaProblema = document.getElementById("dia_problema").value;
    const estado = document.getElementById("estado").value;

    const response = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ empleado, tipo_equipo: tipoEquipo, detalle_problema: detalleProblema, dia_problema: diaProblema, estado })
    });

    if (response.ok) {
        alert("Incidente creado exitosamente");
        fetchIncidents();
    } else {
        alert("Error al crear el incidente");
    }
}

// Actualizar estado de incidente
async function updateIncident() {
    const id = document.getElementById("update-id").value;
    const estado = document.getElementById("update-estado").value;

    const response = await fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ estado })
    });

    if (response.ok) {
        alert("Estado actualizado exitosamente");
        fetchIncidents();
    } else {
        alert("Error al actualizar el estado");
    }
}

// Eliminar incidente
async function deleteIncident() {
    const id = document.getElementById("delete-id").value;

    const response = await fetch(`${API_URL}/${id}`, { method: "DELETE" });

    if (response.ok) {
        alert("Incidente eliminado exitosamente");
        fetchIncidents();
    } else {
        alert("Error al eliminar el incidente");
    }
}
