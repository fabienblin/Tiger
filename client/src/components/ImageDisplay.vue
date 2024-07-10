<script setup lang="js">
import * as d3 from 'd3'

// Get dimensions of the viewport
const width = window.innerWidth;
const height = window.innerHeight;

// Define shapes for different types
const shapeMap = {
    1: d3.symbolCircle,
    2: d3.symbolTriangle,
    3: d3.symbolSquare,
    4: d3.symbolCross
};

// Create SVG element to cover the entire page
const svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height)
    .style("position", "absolute")
    .style("top", 0)
    .style("left", 0)
    .append("g"); // Append a group to properly position elements

// Initialize WebSocket
initWS();

// Tooltip
const tooltip = d3.select("body").append("div")
    .attr("class", "tooltip")
    .style("opacity", 0);

function initWS() {
    // Create a WebSocket connection
    const socket = new WebSocket("ws://127.0.0.1:8080/imageDisplay");

    // Connection is opened
    socket.addEventListener("open", function (event) {
        socket.send("Hello, server!");
    });

    // Listen for messages
    socket.addEventListener("message", function (event) {
        console.log("Message from server:", event.data);
        updateImage(JSON.parse(event.data));
    });

	socket.onclose = function(event) {
		console.log("WebSocket is closed now.");
	};

	socket.onerror = function(error) {
		console.log("WebSocket error observed:", error);
	};

}

function updateImage(graphData) {
    const nodes = Object.values(graphData).map(node => ({
        id: node.id,
        x: node.x * width, // Scale normalized x to viewport width
        y: node.y * height, // Scale normalized y to viewport height
        type: node.type,
        status: node.status
    }));

    const links = Object.values(graphData).flatMap(node =>
        node.neighbours ? Object.values(node.neighbours).map(neighbour => ({
            source: node.id,
            target: neighbour.id
        })) : []
    );

    // Add links
    const link = svg.selectAll("line")
        .data(links, d => `${d.source}-${d.target}`);

    link.enter().append("line")
        .attr("stroke", "#999")
        .attr("stroke-width", "2")
        .merge(link)
        .attr("x1", d => getNodeById(d.source).x)
        .attr("y1", d => getNodeById(d.source).y)
        .attr("x2", d => getNodeById(d.target).x)
        .attr("y2", d => getNodeById(d.target).y);

    link.exit().remove();

    // Add nodes
    const node = svg.selectAll("path")
        .data(nodes, d => d.id);

    node.enter().append("path")
        .attr("class", "node")
        .attr("transform", d => `translate(${d.x},${d.y})`)
        .attr("d", d => d3.symbol()
            .type(shapeMap[d.type])
            .size(500)())
        .merge(node)
        .attr("transform", d => `translate(${d.x},${d.y})`)
        .attr("fill", d => {
            if (d.status.isAlarmed) { return "red"; }
            else if (!d.status.isActive) { return "grey"; }
            else { return "steelblue"; }
        })
        .on("mouseover", handleMouseOver)
        .on("mouseout", handleMouseOut);

    node.exit().remove();

    // Add labels (optional)
    const label = svg.selectAll("text")
        .data(nodes, d => d.id);

    label.enter().append("text")
        .text(d => d.id)
        .attr("dx", 12)
        .attr("dy", 4)
        .style("visibility", "hidden")
        .merge(label)
        .attr("x", d => d.x)
        .attr("y", d => d.y);

    label.exit().remove();

    // Function to handle mouse over event
    function handleMouseOver(event, d) {
        tooltip.transition()
            .duration(200)
            .style("opacity", .9);
        tooltip.html(d.id)
            .style("left", (event.pageX + 10) + "px")
            .style("top", (event.pageY - 20) + "px");

        label.filter(labelData => labelData.id === d.id)
            .style("visibility", "visible");
    }

    // Function to handle mouse out event
    function handleMouseOut() {
        tooltip.transition()
            .duration(500)
            .style("opacity", 0);

        label.style("visibility", "hidden");
    }

    // Helper function to get node by ID
    function getNodeById(id) {
        return nodes.find(node => node.id === id);
    }
}

</script>

<template>
</template>

<style scoped>
</style>
