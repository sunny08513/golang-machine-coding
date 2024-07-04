package main

import (
	"errors"
	"fmt"
)

type Item struct {
	Id       string
	Name     string
	Price    float64
	Quantity int
}

type CartItem struct {
	Item     Item
	Quantity int
}

type Cart struct {
	Items    map[string]*CartItem
	Discount float64
}

type Inventory struct {
	Items map[string]*Item
}

func NewInventory() *Inventory {
	return &Inventory{Items: make(map[string]*Item)}
}

func NewCart() *Cart {
	return &Cart{Items: make(map[string]*CartItem)}
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items[item.Id] = &item
}

func (inv *Inventory) UpdateQuantity(itemID string, quantity int) error {
	item, exists := inv.Items[itemID]
	if !exists {
		return errors.New("item not found")
	}
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	item.Quantity = quantity
	inv.Items[itemID] = item
	return nil
}

func (inv *Inventory) GetItem(itemID string) (*Item, error) {
	item, exists := inv.Items[itemID]
	if !exists {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (c *Cart) AddItem(item Item, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	cartItem, exists := c.Items[item.Id]
	if exists {
		cartItem.Quantity += quantity
	} else {
		cartItem = &CartItem{
			Item:     item,
			Quantity: quantity,
		}
	}
	c.Items[item.Id] = cartItem
	return nil
}

func (c *Cart) RemoveItem(itemID string) error {
	_, exists := c.Items[itemID]
	if !exists {
		return errors.New("item not in cart")
	}
	delete(c.Items, itemID)
	return nil
}

func (c *Cart) ViewCart() {
	for _, cartItem := range c.Items {
		fmt.Printf("Item: %s, Quantity: %d, Price: %.2f\n", cartItem.Item.Name, cartItem.Quantity, cartItem.Item.Price)
	}
	fmt.Printf("Total Price: %.2f\n", c.CalculateTotal())
	if c.Discount > 0 {
		fmt.Printf("Discount: %.2f%%\n", c.Discount)
		fmt.Printf("Discounted Price: %.2f\n", c.CalculateTotalWithDiscount())
	}
}

func (c *Cart) CalculateTotal() float64 {
	total := 0.0
	for _, cartItem := range c.Items {
		total += cartItem.Item.Price * float64(cartItem.Quantity)
	}
	return total
}

func (c *Cart) ApplyDiscount(discount float64) error {
	if discount < 0 || discount > 100 {
		return errors.New("discount must be between 0 and 100")
	}
	c.Discount = discount
	return nil
}

func (c *Cart) CalculateTotalWithDiscount() float64 {
	total := c.CalculateTotal()
	discountAmount := total * (c.Discount / 100)
	return total - discountAmount
}
