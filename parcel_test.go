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
    assert.Equal(t, parcel, parcelFromDB)

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
	assert.Equal(t, newAddress, parcelFromDB.Address)
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
	require.Equal(t, newStatus, parcelFromDB.Status)

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
    // prepare
    db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
    require.NoError(t, err)
    defer db.Close()

    store := NewParcelStore(db)

    parcels := []Parcel{
        getTestParcel(),
        getTestParcel(),
        getTestParcel(),
    }
    parcelMap := map[int]Parcel{}

    // задаём всем посылкам один и тот же идентификатор клиента
    client := randRange.Intn(10_000_000)
    parcels[0].Client = client
    parcels[1].Client = client
    parcels[2].Client = client

    // add
    for i := 0; i < len(parcels); i++ {
        id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
        require.NoError(t, err)
        require.NotEmpty(t, id)
        // обновляем идентификатор добавленной у посылки
        parcels[i].Number = id

        // сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
        parcelMap[id] = parcels[i]
    }

    // get by client
    storedParcels, err := store.GetByClient((client)) // получите список посылок по идентификатору клиента, сохранённого в переменной client
    // убедитесь в отсутствии ошибки
    require.NoError(t, err)
    // убедитесь, что количество полученных посылок совпадает с количеством добавленных
    require.Equal(t, len(parcels), len(storedParcels))
    // check
    for _, parcel := range storedParcels {
        _, ok := parcelMap[parcel.Number]
        // в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
        // убедитесь, что все посылки из storedParcels есть в parcelMap
        assert.True(t, ok)
        // убедитесь, что значения полей полученных посылок заполнены верно
        assert.Equal(t, parcel, parcelMap[parcel.Number])
    }
}
