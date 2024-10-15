import threading
import time
import random

# Semáforos y locks
santa_sleep = threading.Semaphore(0)  # Despierta a Papá Noel cuando hay 3 elfos
elves_waiting_room = threading.Semaphore(3)  # Máximo de 3 elfos en la sala de espera
reindeers_arrival = threading.Semaphore(0)  # Avisa a Papá Noel cuando llegan todos los renos
mutex = threading.Lock()  # Control de acceso a los contadores de elfos y renos
elf_count_lock = threading.Lock()
reindeer_count_lock = threading.Lock()

# Variables compartidas
elf_count = 0
reindeer_count = 0
elves_finished = 0

NUM_JOGUINES = 3
NUM_ELFS = 6
CAPACITAT_SALA_ESPERA = 3
NUM_RENS = 9

class SantaClaus(threading.Thread):
    def run(self):
        global elf_count, reindeer_count
        print("-------> El Pare Noel diu: Estic despert però me'n torn a jeure")# Comença la simulació
        while True:
            santa_sleep.acquire()  # Espera a ser despertado por el tercer elfo o el último reno
            if elf_count == 3:
                print("-------> El Pare Noel diu: Atendré els dubtes d'aquests 3")# El tercer elf l'ha despert
                time.sleep(2)  # Simula el tiempo de respuesta
                print("-------> El Pare Noel diu: Estic cansat me'n torn a jeure")# Ha resolt els tres dubtes
                # Permitir a más elfos que tengan dudas entrar en la sala de espera
                elf_count_lock.acquire()
                elf_count = 0  # Resetea el conteo de elfos que hicieron preguntas
                elf_count_lock.release()

            if elves_finished == 6:
                print("-------> Pare Noel diu: Les joguines estan llestes. I Els rens?")# Els elfs han acabat

            if elves_finished == 6 and reindeer_count == 9:
                print("-------> Pare Noel diu: Enganxaré els rens i partiré")# Tots els rens han arribat
                for _ in range(9):  # Engancha cada reno al trineo
                    reindeers_arrival.release()
                print("-------> El Pare Noel ha enganxat els rens, ha carregat les joguines i se'n va")# Ha acabat
                break

class Elf(threading.Thread):
    def __init__(self, name):
        super().__init__()
        self.name = name

    def run(self):
        global elf_count, elves_finished
        print(f"Hola som l'elf {self.name} construiré {NUM_JOGUINES} joguines")# Comença la simulació
        for i in range(1, NUM_JOGUINES):
            time.sleep(random.randint(1, 3))  # Simula el tiempo de construcción
            print(f"{self.name} diu: tinc dubtes amb la joguina {i}")# Entra a la sala de espera
            elves_waiting_room.acquire()
            
            with elf_count_lock:
                elf_count += 1
                if elf_count == 3:
                    print(f"{self.name} diu: Som 3 que tenim dubtes, PARE NOEEEEEL!")# El tercer elf desperta al Pare Noel
                    santa_sleep.release()  # Despierta a Papá Noel
            
            time.sleep(1)  # Simula la espera de ayuda
            print(f"{self.name} diu: Construeixo la joguina amb ajuda")# Ha aclarit el dubte i surt de la sala de espera
            elves_waiting_room.release()
        
        with elf_count_lock:
            elves_finished += 1
            if elves_finished == NUM_ELFS:
                print(f"{self.name} diu: Som el darrer avisaré al Pare Noel")# El darrer elf ha acabat la seva tercera joguina
                santa_sleep.release()  # Último elfo avisa a Papá Noel

        print(f"L'elf {self.name} ha fet les seves joguines i acaba <---------")# Ha acabat

class Reindeer(threading.Thread):
    def __init__(self, name):
        super().__init__()
        self.name = name

    def run(self):
        global reindeer_count
        print(f"{self.name} se'n va a pasturar")# Comença la simulació
        time.sleep(random.randint(5,9))  # Simula el tiempo de llegada del reno

        with reindeer_count_lock:
            reindeer_count += 1
            if reindeer_count == NUM_RENS:
                print(f"El ren {self.name} diu: Som el darrer en voler podem partir")# Ha arribat el darrer ren
                santa_sleep.release()  # Despierta a Papá Noel cuando todos los renos están listos
            else:
                print(f"El ren {self.name} arriba, {reindeer_count}")# Ha arribat un ren que no és el darrer

        reindeers_arrival.acquire()
        print(f"El ren {self.name} està enganxat al trineu")# El ren ha estat enganxat i acaba

# Inicializar y ejecutar los hilos
santa = SantaClaus()
santa.start()

elves = [Elf(f"Elf-{i + 1}") for i in range(NUM_ELFS)]
for elf in elves:
    elf.start()

reindeers = [Reindeer(f"Reindeer-{i + 1}") for i in range(NUM_RENS)]
for reindeer in reindeers:
    reindeer.start()

# Esperar a que todos los hilos terminen
for elf in elves:
    elf.join()
for reindeer in reindeers:
    reindeer.join()
santa.join()
