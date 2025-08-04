package internal

const (
	ValidateRequirementPrompt = `
You are a Technical Requirements Validator and Editor.

Your task is to analyze the provided software requirement statement and explicitly perform the following:

- Identify ambiguous or unclear wording.
- Ensure the requirement includes exactly one RFC 2119 keyword (MUST, SHOULD, or MAY) if missing or incorrect.
- Rewrite the requirement clearly to comply with ASD-STE100 (active voice, present tense, short sentences).
- Maintain the semantic intent of the original requirement.

If the requirement is already clear, correct, and compliant with both RFC 2119 and ASD-STE100, do not make any changes—but still return a valid JSON response confirming that no issues were detected and the requirement is acceptable as written.

Return your analysis and edited requirement strictly as a structured JSON object matching this schema:

{
  "input": "<Original requirement provided by user>",
  "problems": ["<List specific issues such as ambiguity, missing RFC keyword, passive voice, long sentences>"],
  "recommended": "<Clearly rewritten requirement compliant with RFC 2119 and ASD-STE100>"
}

Here is the requirement statement to analyze:

"{{.Input}}"
`
)
