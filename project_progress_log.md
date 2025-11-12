# Project Progress Log (2025-11-07)

## Summary of Work Done

1.  **Fixed the login issue:**
    *   Diagnosed and confirmed the absence of a default user in the database.
    *   Enabled the `dev` build tag to activate the creation of a temporary user (`admin`/`password123`), allowing the login process to function correctly.

2.  **Implemented Frontend Data Visualization:**
    *   **Enriched the vehicle list:** Corrected the data model to display the model, battery level, and real-time status for each vehicle in the list.
    *   **Created a vehicle detail page:** Implemented functionality to navigate to a detail page by clicking on a vehicle. This page displays all historical decision logs for the vehicle as cards, including images, actions, and confidence scores.
    *   **Added a global statistics chart:** Introduced a new backend API to fetch all decision logs and added a pie chart to the main dashboard to visualize the distribution of decision actions (e.g., "pickup", "ignore").

3.  **Fixed Frontend Compilation and Memory Leak Issues:**
    *   Resolved Material-UI `Grid` component compilation errors by refactoring the layout to use a more robust CSS Grid implementation with the `Box` component.
    *   Fixed a memory leak and page crash caused by the chart component by providing it with a fixed-height container, preventing infinite resizing.

---

## API Interface Changes

### New Endpoints

#### `GET /api/v1/decision-logs`

-   **Description:** Fetches all decision logs from all vehicles.
-   **Authentication:** Required (JWT Bearer Token).
-   **Response Body:**

    ```json
    {
      "logs": [
        {
          "id": "...",
          "vehicle_id": "...",
          "timestamp": "...",
          "image_url": "...",
          "server_decision": {
            "image_id": "...",
            "action": "pickup",
            "confidence": 0.95,
            "reason": "is_trash_type_A"
          }
        }
      ],
      "total": 1
    }
    ```

### Modified Endpoints

#### `GET /api/v1/vehicles/:id/decision-logs`

-   **Description:** The response format for this endpoint has been updated for consistency.
-   **Old Response Format:**
    ```json
    {
      "data": [...],
      "pagination": {...}
    }
    ```
-   **New Response Format:**
    ```json
    {
      "logs": [...],
      "total": ...
    }
    ```

---

## Next Steps

1.  **Integrate ONNX Model (High Challenge):**
    *   **Goal:** Enable the backend service to load and run the provided `.onnx` model for real-time AI image recognition, replacing the current mocked data.
    *   **Steps:** This requires installing the ONNX Runtime C library in your development environment and writing the necessary CGo code to interface with it. This is a complex but core-functional task.

2.  **Integrate `qwen3-vl` Vision-Language Model API:**
    *   **Goal:** Add advanced image understanding capabilities to the system, such as generating natural language descriptions of images or answering specific questions about them.
    *   **Steps:** This involves creating a new service in the backend to call the `qwen3-vl` API and integrating the results into our application.
