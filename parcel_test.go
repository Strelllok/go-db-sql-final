package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
    db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

    store := NewParcelStore(db)
    parcel := getTestParcel()

    id, err := store.Add(parcel)
    require.NoError(t, err)
    require.NotEmpty(t, id)

    parcelFromDB, err := store.Get(id)
    require.NoError(t, err)
	parcel.Number = id
    require.Equal(t, parcel, parcelFromDB)

    err = store.Delete(id)
    require.NoError(t, err)

    _, err = store.Get(id)
    assert.Error(t, err)
	assert.Equal(t, parcel, parcelFromDB)

}
// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	parcelFromDB, err := store.Get(id)
	assert.NoError(t, err)
	assert.NotEqual(t, newAddress, parcelFromDB.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newStatus := ParcelStatusDelivered
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	parcelFromDB, err := store.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, newStatus, parcelFromDB.Status)

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
    db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

    store := NewParcelStore(db)
    parcels := []Parcel{
        getTestParcel(),
        getTestParcel(),
        getTestParcel(),
    }
    parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
    parcels[0].Client = client
    parcels[1].Client = client
    parcels[2].Client = client

    for i := 0; i < len(parcels); i++ {
        id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
        require.NoError(t, err)
        require.NotEmpty(t, id)
        // обновляем идентификатор добавленной у посылки
        parcels[i].Number = id

        // сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
        parcelMap[id] = parcels[i]
	}

    storedParcels, err := store.GetByClient(client)

    require.NoError(t, err)
	require.Len(t, len(parcels), len(storedParcels))

    for _, parcel := range storedParcels {
        expectedParcel, ok := parcelMap[parcel.Number]
        assert.True(t, ok)
        assert.Equal(t, parcel, expectedParcel)
    }
}
