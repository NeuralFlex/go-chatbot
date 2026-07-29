package services

const PresetCostSavings = `Analyse the ledger data below and identify cost-saving opportunities. For each opportunity:
- Name the account and current amount
- State the saving potential in currency (realistic estimate)
- Give one concrete action the controller can take immediately

Focus only on costs that are unusually high relative to revenue or that have grown more than 15% vs the comparison period. Ignore fixed costs that cannot be reduced. Answer in bullet points, max 800 characters.`

const PresetMeetingNotes = `You are preparing for a client meeting today.
Based on the ledger data below, give a concise pre-meeting briefing.

Structure your answer in exactly four parts:

- KEY NUMBERS: the 3 most important figures the client will ask about
  (revenue, margin or cash, with EUR amounts and % change vs comparison period)

- POSITIVE DEVELOPMENTS: max 2 items the controller can highlight confidently

- CONCERNS TO ADDRESS: max 2 items the client is likely to raise.
  For each, include a brief factual explanation the controller can give

- ONE OPEN QUESTION: the single most important thing that needs
  clarification or a decision from the client at this meeting

Max 800 characters.`

const PresetClientProposals = `You are preparing recommendations for a client meeting.
Based on the ledger data below, identify what to propose to the client at this meeting.

Structure your answer in exactly three parts:

- QUICK WINS (this month):
  Max 2 concrete proposals the client can act on immediately.
  For each: what to do, expected impact in EUR, and why now.

- STRATEGIC PROPOSALS (next 1 to 3 months):
  Max 2 proposals that require planning or investment.
  For each: what to propose, business rationale in one sentence,
  and what data or decision is needed from the client.

- ONE RISK TO RAISE:
  The single most important risk visible in the ledger that the client
  may not be aware of. State the risk, the EUR exposure, and
  one action that would reduce it.

Max 800 characters.`
