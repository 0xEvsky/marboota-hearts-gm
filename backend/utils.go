package main

type Queue struct {
	items []*Client
}

func (q *Queue) Enqueue(c *Client) {
	q.items = append(q.items, c)
}

func (q *Queue) Dequeue() (*Client, bool) {
	if len(q.items) == 0 {
		return nil, false
	}
	client := q.items[0]
	q.items = q.items[1:]
	return client, true
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Peek() (*Client, bool) {
	if len(q.items) == 0 {
		return nil, false
	}
	return q.items[0], true
}
