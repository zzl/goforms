package main

import "github.com/zzl/goforms/framework/virtual"

// Animal is the public interface client code faces.
type Animal interface {
	virtual.Virtual
	Greet()
}

// AnimalSpi is the protected interface used for implementation.
type AnimalSpi interface {
	Say(something string)
}

// AnimalInterface is the full interface that combines Animal and AnimalSpi
type AnimalInterface interface {
	Animal
	AnimalSpi
}

// AnimalObject implements AnimalInterface.
type AnimalObject struct {

	// VirtualObject should be embedded as the parent object
	// for types at the root of the type hierarchy
	virtual.VirtualObject[AnimalInterface]
}

// Greet implements Animal.Greet
func (this *AnimalObject) Greet() {
	this.RealObject.Say("Hello!")
}

// Say implements AnimalSpi.Say
func (this *AnimalObject) Say(something string) {
	//abstract method, do nothing
}

// Cat is derived from Animal
type Cat interface {
	Animal
}

type CatSpi interface {
	AnimalSpi
	GetVoice() string
}

type CatInterface interface {
	Cat
	CatSpi
}

type CatObject struct {
	AnimalObject
}

func (this *CatObject) GetVoice() string {
	return "Meow"
}

func (this *CatObject) Say(something string) {
	spi := this.RealObject.(CatSpi)
	println(spi.GetVoice() + "! " + something)
}

// Dog is derived from Animal
type Dog interface {
	Animal
}

type DogSpi interface {
	AnimalSpi
}

type DogInterface interface {
	Dog
	DogSpi
}

type DogObject struct {
	AnimalObject
}

func (this *DogObject) Say(something string) {
	println("Woof! " + something)
}

// Tiger is derived from Cat
type Tiger interface {
	Cat
}

type TigerSpi interface {
	CatSpi
}

type TigerInterface interface {
	Tiger
	TigerSpi
}

type TigerObject struct {
	CatObject
	super *CatObject
}

func (this *TigerObject) GetVoice() string {
	return "Roar~~"
}

func (this *TigerObject) Greet() {
	this.super.Greet()
	this.Say("Welcome to the tiger world!")
}

// main
func main() {
	var animals []Animal

	var cat Cat = virtual.New[CatObject]()
	var dog Dog = virtual.New[DogObject]()
	var tiger Tiger = virtual.New[TigerObject]()

	animals = append(animals, cat)
	animals = append(animals, dog)
	animals = append(animals, tiger)

	for _, animal := range animals {
		animal.Greet()
	}

	/*expected output:

	  Meow! Hello!
	  Woof! Hello!
	  Roar~~! Hello!
	  Roar~~! Welcome to the tiger world!
	*/
}
