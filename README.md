# goldsteps
Task Management and Financial Information Collection
![dashboard_20250328](https://github.com/user-attachments/assets/4c8d3ddc-1082-4d87-889b-3bd8b53fca29)

## Tech Stat
* Backend : Go 
* Server-side Framework : Echo

* Frontend : Node.js (npm)
* Client-side Framework : Next.js

* Docker

## Setup
To execute services, run `cd goldsteps && make local` for local or `cd goldsteps && make build && make up` for Docker.

To store master data of stocks, run `curl http://localhost:8080/api/stock_master`.

## How to Use
### Task Management
Dashboard
<img width="1512" alt="task_management" src="https://github.com/user-attachments/assets/9e34d42a-ee87-48df-916f-64bdc1375d88" />

Create or Edit new tasks in the modal

<img width="350" alt="create_event" src="https://github.com/user-attachments/assets/64d246de-96f7-42aa-90a3-8eba3d87e6db" />

### Financial Information Collection
Tool list
<img width="302" alt="news_tool" src="https://github.com/user-attachments/assets/832df3a6-0ae0-4f98-8ce1-77c9f3272ae7" />

Press `View Headline`, then it crawls the latest news in [Bloomberg](https://www.bloomberg.co.jp/).

<img width="1511" alt="headline_zone" src="https://github.com/user-attachments/assets/2dc30958-d961-48f9-bcd8-979da727d836" />

Press `Add to Milestone` in an article, then it adds the news into the milestones list.

<img width="1512" alt="add_to_milestone" src="https://github.com/user-attachments/assets/aee3367b-422a-4f31-aa29-839ab9402cd2" />

Press `Show Milestones`, then the list is shown.

<img width="450" alt="milestone_list" src="https://github.com/user-attachments/assets/3ebc5a84-4e1f-43ee-b4be-0ab674d2345a" />

---

In Stocks Search page, collecting data and news of the codes

<img width="1117" alt="stock_search_loading" src="https://github.com/user-attachments/assets/d74eccf9-feb2-4260-b515-d5db4241e35f" />

As the result;

<img width="1147" alt="stock_search" src="https://github.com/user-attachments/assets/bbecf590-bdb4-40ab-b460-19a44ff98401" />

The history of inputted codes is cached.

<img width="1135" alt="stock_search_history" src="https://github.com/user-attachments/assets/f2aa4d05-77a0-4161-9da8-fde2f224a20b" />

---

In Arts Search page, exploring the previous news related to the inputted free word

<img width="455" alt="article_search" src="https://github.com/user-attachments/assets/992ae804-c1b8-4e7c-99b2-73340b0c5b7f" />

The history of free words is cached.

<img width="455" alt="arts_search_history" src="https://github.com/user-attachments/assets/933383b8-2388-4052-9623-e8377d9facfc" />

