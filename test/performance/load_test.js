import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

export const options = {
  scenarios: {
    requirement_check: {
      executor: "constant-arrival-rate",
      rate: 5000,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 100,
      maxVUs: 1000,
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<3"],
    errors: ["rate<0.01"],
  },
  // Add minimal output configuration
  systemTags: ["status", "method", "url"],
  summaryTrendStats: ["avg", "min", "med", "max", "p(90)", "p(95)"],
};

const errorRate = new Rate("errors");

function generateTestData(questionCount) {
  const data = [];
  for (let i = 0; i < questionCount; i++) {
    data.push({
      question_id: Math.floor(Math.random() * 100) + 300,
      answer_id: Math.floor(Math.random() * 100) + 500,
    });
  }
  return data;
}

export default function() {
  const payload = {
    country: "US",
    data: generateTestData(Math.floor(Math.random() * 10) + 1),
  };

  const response = http.post(
    "http://localhost:8080/transform",
    JSON.stringify(payload),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  const success = check(response, {
    "status is 200": (r) => r.status === 200,
    "response has data": (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && Array.isArray(body.data);
      } catch {
        return false;
      }
    },
    "response time < 3ms": (r) => r.timings.duration < 3,
  });

  errorRate.add(!success);
  sleep(0.1);
}
