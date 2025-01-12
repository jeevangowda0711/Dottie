// frontend/pages/index.js (Next.js example)
"use client";

import { useState } from "react";

export default function Home() {
  const [symptoms, setSymptoms] = useState<string[]>([]);
  const [cycleLength, setCycleLength] = useState(0);
  const [cycleDuration, setCycleDuration] = useState(0);
  const [age, setAge] = useState(0);
  const [analysisResult, setAnalysisResult] = useState<AnalysisResult | null>(null);

  interface AnalysisResult {
    is_normal: boolean;
    abnormalities?: string[];
    recommendations: string[];
  }

  interface AnalyzeSymptomsResponse {
    symptoms: string[];
    cycle_length: number;
    cycle_duration: number;
    age: number;
    flow: string;
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const res = await fetch("http://localhost:8080/analyzeSymptoms", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        symptoms,
        cycle_length: cycleLength,
        cycle_duration: cycleDuration,
        age,
        flow: "heavy",
      } as AnalyzeSymptomsResponse),
    });
    const data: AnalysisResult = await res.json();
    setAnalysisResult(data);
  };

  return (
    <div>
      <h1>Menstrual Symptom Checker</h1>
      <form onSubmit={handleSubmit}>
        <label>Cycle Length (days):
          <input
            type="number"
            value={cycleLength}
            onChange={(e) => setCycleLength(+e.target.value)}
          />
        </label>
        <label>Cycle Duration (days):
          <input
            type="number"
            value={cycleDuration}
            onChange={(e) => setCycleDuration(+e.target.value)}
          />
        </label>
        <label>Age:
          <input
            type="number"
            value={age}
            onChange={(e) => setAge(+e.target.value)}
          />
        </label>
        <label>Symptoms (comma-separated):
          <input
            type="text"
            onChange={(e) => setSymptoms(e.target.value.split(","))}
          />
        </label>
        <button type="submit">Analyze</button>
      </form>

      {analysisResult && (
        <div>
          <h2>Analysis Result</h2>
          <p>IsNormal: {analysisResult.is_normal ? "Yes" : "No"}</p>
          {analysisResult.abnormalities && (
            <p>Abnormalities: {analysisResult.abnormalities.join(", ")}</p>
          )}
          <p>Recommendations: {analysisResult.recommendations.join(", ")}</p>
        </div>
      )}
    </div>
  );
}
