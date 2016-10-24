#!/usr/bin/env ruby

ROOT_PATH = 'assets'
TGFF_PATH = 'tgff'
TGFF_BASE = 'base.tgffopt'
QUANTITIES = ['delay', 'energy', 'temperature']

class Application
  def initialize(cores, tasks)
    @cores, @tasks = cores, tasks
    @base = '%03d_%03d' % [cores, tasks]
  end

  def run(quantities = QUANTITIES)
    self.tgff
    quantities.each do |quantity|
      self.json(quantity)
    end
  end

  def tgff(input = TGFF_BASE)
    raise("'#{input}' does not exist") unless File.exist?(input)

    template = Template.new(input)

    output = File.join("#{@base}.tgffopt")
    sample = File.join("#{@base}.tgff")

    puts("Generating #{@base}...")

    seed = @cores * @tasks
    loop do
      File.open(output, 'w') do |file|
        template.complete(seed: seed, cores: @cores, tasks: @tasks) do |line|
          file << line
        end
      end

      command = "#{TGFF_PATH} #{@base}"
      raise("failed to execute '#{command}'") unless system(command)

      found = File.open(sample).readlines.count do |line|
        line =~ /\s*TASK\s+/
      end
      break if found == @tasks

      seed += 1
    end
  end

  def json(quantity)
    output = "#{@base}_#{quantity}.json"
    inherit = File.join(ROOT_PATH, "%03d_#{quantity}.json" % @cores)
    specification = File.join(ROOT_PATH, "#{@base}.tgff")

    File.open(output, 'w') do |file|
      file << <<-CONTENT
{
	"inherit": "#{inherit}",

	"system": {
		"specification": "#{specification}"
	}
}
      CONTENT
    end
  end
end

class Template
  def initialize(filename)
    @template = File.open(filename).readlines
  end

  def complete(options)
    parameters = {}
    options.each do |name, value|
      parameters["[#{name.to_s.upcase}]"] = value.to_s
    end

    @template.each do |line|
      parameters.each do |name, value|
        line = line.gsub(name, value)
      end
      yield(line)
    end
  end
end

if ARGV.length != 2
  puts 'Usage: generate <number of cores> <number of tasks>'
  exit
end

Application.new(ARGV[0].to_i, ARGV[1].to_i).run
